package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"productservice/internal/domain"
)

type PostgresTagRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewPostgresTagRepository(db *sql.DB, logger *slog.Logger) *PostgresTagRepository {
	return &PostgresTagRepository{
		db:     db,
		logger: logger,
	}
}

// Create a new tag for a shop
func (r *PostgresTagRepository) Create(ctx context.Context, tag domain.Tag) (*domain.Tag, error) {
	query := `
		INSERT INTO tags (shop_id, name, slug)
		VALUES ($1, $2, $3)
		RETURNING id, shop_id, name, slug, created_at, updated_at
	`
	err := r.db.QueryRowContext(ctx, query, tag.ShopID, tag.Name, tag.Slug).Scan(
		&tag.ID, &tag.ShopID, &tag.Name, &tag.Slug, &tag.CreatedAt, &tag.UpdatedAt,
	)

	if err != nil {
		r.logger.ErrorContext(ctx, "failed to create tag", "error", err, "name", tag.Name)
		return nil, err
	}
	return &tag, nil
}

// GetByID fetches a single tag, optionally with the product count
func (r *PostgresTagRepository) GetByID(ctx context.Context, shopID, id string, includeCount bool) (*domain.Tag, int32, error) {
	countSubquery := "0 as product_count"
	if includeCount {
		countSubquery = "(SELECT count(1) FROM product_tags pt WHERE pt.tag_id = tags.id) as product_count"
	}

	query := fmt.Sprintf(`
		SELECT id, shop_id, name, slug, created_at, updated_at, %s
		FROM tags
		WHERE id = $1 AND shop_id = $2 AND deleted_at IS NULL
	`, countSubquery)

	var t domain.Tag
	var count int32
	err := r.db.QueryRowContext(ctx, query, id, shopID).Scan(
		&t.ID, &t.ShopID, &t.Name, &t.Slug, &t.CreatedAt, &t.UpdatedAt, &count,
	)

	if err != nil {
		return nil, 0, err
	}
	return &t, count, nil
}

// GetDetail fetches a tag and a paginated list of products associated with it
func (r *PostgresTagRepository) GetDetail(ctx context.Context, shopID, id string, page, pageSize int) (*domain.Tag, []*domain.ProductSummary, int32, error) {
	if pageSize <= 0 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	// 1. Get the Tag
	tag, _, err := r.GetByID(ctx, shopID, id, false)
	if err != nil {
		return nil, nil, 0, err
	}

	// 2. Get Total Count of products for this tag
	var total int32
	err = r.db.QueryRowContext(ctx, `
		SELECT count(1) FROM product_tags WHERE tag_id = $1
	`, id).Scan(&total)
	if err != nil {
		return nil, nil, 0, err
	}

	// 3. Get associated products (Lightweight Summary)
	query := `
		SELECT p.id, p.name, p.sku, p.price, p.is_active
		FROM products p
		JOIN product_tags pt ON pt.product_id = p.id
		WHERE pt.tag_id = $1 AND p.deleted_at IS NULL
		ORDER BY p.created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, query, id, pageSize, offset)
	if err != nil {
		return nil, nil, 0, err
	}
	defer rows.Close()

	var products []*domain.ProductSummary
	for rows.Next() {
		var p domain.ProductSummary
		var sku sql.NullString
		if err := rows.Scan(&p.ID, &p.Name, &sku, &p.Price, &p.IsActive); err != nil {
			return nil, nil, 0, err
		}
		if sku.Valid {
			p.SKU = sku.String
		}
		products = append(products, &p)
	}

	return tag, products, total, nil
}

// List tags for a shop with optional search and product counts
func (r *PostgresTagRepository) List(
	ctx context.Context,
	shopID string,
	search string,
	includeCount bool,
	page, pageSize int,
) ([]*domain.TagWithCount, int64, error) {

	if pageSize <= 0 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	baseWhere := "WHERE shop_id = $1 AND deleted_at IS NULL"
	args := []any{shopID}
	argPos := 2

	if search != "" {
		baseWhere += fmt.Sprintf(" AND name ILIKE $%d", argPos)
		args = append(args, "%"+search+"%")
		argPos++
	}

	// Get Total Count
	var total int64
	err := r.db.QueryRowContext(ctx, "SELECT count(1) FROM tags "+baseWhere, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Main Query
	countSubquery := "0 as product_count"
	if includeCount {
		countSubquery = "(SELECT count(1) FROM product_tags pt WHERE pt.tag_id = tags.id) as product_count"
	}

	query := fmt.Sprintf(`
		SELECT id, shop_id, name, slug, created_at, updated_at, %s
		FROM tags
		%s
		ORDER BY name ASC
		LIMIT $%d OFFSET $%d
	`, countSubquery, baseWhere, argPos, argPos+1)

	args = append(args, pageSize, offset)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var tags []*domain.TagWithCount
	for rows.Next() {
		var t domain.TagWithCount
		err := rows.Scan(&t.ID, &t.ShopID, &t.Name, &t.Slug, &t.CreatedAt, &t.UpdatedAt, &t.ProductCount)
		if err != nil {
			return nil, 0, err
		}
		tags = append(tags, &t)
	}

	return tags, total, nil
}

// AssignTagsToProduct manages the many-to-many relationship
func (r *PostgresTagRepository) AssignToProduct(
	ctx context.Context,
	shopID, productID string,
	tagIDs []string,
	replace bool,
) ([]*domain.Tag, error) {

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if replace {
		_, err = tx.ExecContext(ctx, "DELETE FROM product_tags WHERE product_id = $1", productID)
		if err != nil {
			return nil, err
		}
	}

	for _, tagID := range tagIDs {
		// Ensure tag belongs to the shop before assigning
		_, err = tx.ExecContext(ctx, `
			INSERT INTO product_tags (product_id, tag_id)
			SELECT $1, $2 WHERE EXISTS (SELECT 1 FROM tags WHERE id = $2 AND shop_id = $3)
			ON CONFLICT DO NOTHING
		`, productID, tagID, shopID)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return r.GetByProductID(ctx, shopID, productID)
}

func (r *PostgresTagRepository) RemoveFromProduct(
	ctx context.Context,
	shopID string,
	productID string,
	tagIDs []string,
) error {
	if len(tagIDs) == 0 {
		return nil
	}

	r.logger.DebugContext(ctx, "removing tags from product",
		"productID", productID,
		"shopID", shopID,
		"tagCount", len(tagIDs),
	)

	// We use a query that ensures we only delete associations where the
	// tag actually belongs to the shop provided in the context.
	query := `
		DELETE FROM product_tags
		WHERE product_id = $1
		  AND tag_id IN (
			  SELECT id FROM tags 
			  WHERE id = ANY($2) AND shop_id = $3
		  )
	`

	res, err := r.db.ExecContext(ctx, query, productID, tagIDs, shopID)
	if err != nil {
		r.logger.ErrorContext(ctx, "failed to remove tags from product",
			"error", err,
			"productID", productID,
			"shopID", shopID,
		)
		return err
	}

	rows, _ := res.RowsAffected()
	r.logger.InfoContext(ctx, "tags removed successfully",
		"productID", productID,
		"rowsAffected", rows,
	)

	return nil
}

// GetByProductID returns all tags associated with a specific product
func (r *PostgresTagRepository) GetByProductID(ctx context.Context, shopID, productID string) ([]*domain.Tag, error) {
	query := `
		SELECT t.id, t.shop_id, t.name, t.slug, t.created_at, t.updated_at
		FROM tags t
		JOIN product_tags pt ON pt.tag_id = t.id
		WHERE pt.product_id = $1 AND t.shop_id = $2 AND t.deleted_at IS NULL
	`
	rows, err := r.db.QueryContext(ctx, query, productID, shopID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*domain.Tag
	for rows.Next() {
		var t domain.Tag
		if err := rows.Scan(&t.ID, &t.ShopID, &t.Name, &t.Slug, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tags = append(tags, &t)
	}
	return tags, nil
}

// Delete a tag (Soft Delete)
func (r *PostgresTagRepository) Delete(ctx context.Context, shopID, id string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE tags SET deleted_at = now() 
		WHERE id = $1 AND shop_id = $2 AND deleted_at IS NULL
	`, id, shopID)
	return err
}
