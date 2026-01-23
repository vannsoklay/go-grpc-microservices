package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	pagination "hpkg/constants"
	sharedErr "hpkg/errors"
	"productservice/domain"
)

type PostgresProductRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewPostgresProductRepository(db *sql.DB, logger *slog.Logger) *PostgresProductRepository {
	return &PostgresProductRepository{
		db:     db,
		logger: logger,
	}
}

func (r *PostgresProductRepository) ListByShopID(
	ctx context.Context,
	shopID,
	search,
	filter,
	sortColumn string,
	sortDesc bool,
	limit int,
	cursor string,
) ([]*domain.Product, string, error) {

	if limit <= 0 || limit > 50 {
		limit = 20
	}

	args := []any{shopID}
	argPos := 2

	query := `
		SELECT id, shop_id, owner_id, name, category, price, detail, created_at, updated_at, deleted_at
		FROM products
		WHERE shop_id = $1
		  AND deleted_at IS NULL
	`

	// Search by name
	if search != "" {
		query += fmt.Sprintf(" AND name ILIKE $%d", argPos)
		args = append(args, "%"+search+"%")
		argPos++
	}

	// Filter by category
	if filter != "" {
		query += fmt.Sprintf(" AND LOWER(category) = LOWER($%d)", argPos)
		args = append(args, filter)
		argPos++
	}

	// Cursor pagination
	if cursor != "" {
		c, err := pagination.DecodeCursor(cursor)
		if err != nil {
			r.logger.ErrorContext(ctx, "failed to decode cursor",
				"error", err,
				"cursor", cursor,
				"shopID", shopID,
			)
			return nil, "", err
		}
		query += fmt.Sprintf(" AND (created_at, id) > ($%d, $%d)", argPos, argPos+1)
		args = append(args, c.CreatedAt, c.ID)
		argPos += 2
	}

	order := "ASC"
	if sortDesc {
		order = "DESC"
	}
	query += fmt.Sprintf(" ORDER BY %s %s, id ASC LIMIT $%d", sortColumn, order, argPos)
	args = append(args, limit+1)

	r.logger.DebugContext(ctx, "executing list query",
		"shopID", shopID,
		"search", search,
		"filter", filter,
		"sortColumn", sortColumn,
		"sortDesc", sortDesc,
		"limit", limit,
	)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.logger.ErrorContext(ctx, "failed to query products",
			"error", err,
			"shopID", shopID,
			"query", query,
			"args", args,
		)
		return nil, "", err
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(
			&p.ID, &p.ShopID, &p.OwnerID,
			&p.Name, &p.Category, &p.Price, &p.Detail,
			&p.CreatedAt, &p.UpdatedAt, &p.DeletedAt,
		); err != nil {
			r.logger.ErrorContext(ctx, "failed to scan product row",
				"error", err,
				"shopID", shopID,
			)
			return nil, "", err
		}
		products = append(products, &p)
	}

	if err := rows.Err(); err != nil {
		r.logger.ErrorContext(ctx, "error iterating product rows",
			"error", err,
			"shopID", shopID,
			"productsCount", len(products),
		)
		return nil, "", err
	}

	r.logger.DebugContext(ctx, "products fetched successfully",
		"shopID", shopID,
		"count", len(products),
	)

	var nextCursor string
	if len(products) > limit {
		last := products[limit-1]
		c, err := pagination.EncodeCursor(pagination.ProductCursor{
			CreatedAt: last.CreatedAt,
			ID:        last.ID,
		})
		if err != nil {
			r.logger.ErrorContext(ctx, "failed to encode cursor",
				"error", err,
				"productID", last.ID,
				"shopID", shopID,
			)
			return nil, "", err
		}
		nextCursor = c
		products = products[:limit]
	}

	return products, nextCursor, nil
}

func (r *PostgresProductRepository) GetByID(
	ctx context.Context,
	id string,
) (*domain.Product, error) {

	r.logger.DebugContext(ctx, "fetching product by id",
		"productID", id,
	)

	row := r.db.QueryRowContext(ctx, `
		SELECT
			id,
			shop_id,
			owner_id,
			name,
			category,
			price,
			detail,
			created_at,
			updated_at,
			deleted_at
		FROM products
		WHERE id = $1
		  AND deleted_at IS NULL
	`, id)

	var p domain.Product
	if err := row.Scan(
		&p.ID,
		&p.ShopID,
		&p.OwnerID,
		&p.Name,
		&p.Category,
		&p.Price,
		&p.Detail,
		&p.CreatedAt,
		&p.UpdatedAt,
		&p.DeletedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			r.logger.WarnContext(ctx, "product not found",
				"productID", id,
			)
			return nil, sharedErr.ErrNotFound
		}
		r.logger.ErrorContext(ctx, "failed to scan product",
			"error", err,
			"productID", id,
		)
		return nil, err
	}

	r.logger.DebugContext(ctx, "product fetched successfully",
		"productID", id,
		"shopID", p.ShopID,
	)

	return &p, nil
}

func (r *PostgresProductRepository) Create(
	ctx context.Context,
	req domain.CreateProductRequest,
) (*domain.Product, error) {

	r.logger.DebugContext(ctx, "creating product",
		"shopID", req.ShopID,
		"ownerID", req.OwnerID,
		"name", req.Name,
		"category", req.Category,
		"price", req.Price,
	)

	row := r.db.QueryRowContext(ctx, `
		INSERT INTO products (
			shop_id,
			owner_id,
			name,
			category,
			description,
			price,
			detail
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING
			id,
			shop_id,
			owner_id,
			name,
			category,
			price,
			detail,
			created_at,
			updated_at
	`,
		req.ShopID,
		req.OwnerID,
		req.Name,
		req.Description,
		req.Category,
		req.Price,
		req.Detail,
	)

	var p domain.Product
	if err := row.Scan(
		&p.ID,
		&p.ShopID,
		&p.OwnerID,
		&p.Name,
		&p.Category,
		&p.Price,
		&p.Detail,
		&p.CreatedAt,
		&p.UpdatedAt,
	); err != nil {
		r.logger.ErrorContext(ctx, "failed to create product",
			"error", err,
			"shopID", req.ShopID,
			"ownerID", req.OwnerID,
			"name", req.Name,
		)
		return nil, err
	}

	r.logger.InfoContext(ctx, "product created successfully",
		"productID", p.ID,
		"shopID", p.ShopID,
	)

	return &p, nil
}

func (r *PostgresProductRepository) Update(
	ctx context.Context,
	req domain.UpdateProductRequest,
) (*domain.Product, error) {

	r.logger.DebugContext(ctx, "updating product",
		"productID", req.ID,
		"name", req.Name,
		"category", req.Category,
		"price", req.Price,
	)

	row := r.db.QueryRowContext(ctx, `
		UPDATE products
		SET
			name = $2,
			category = $3,
			price = $4,
			detail = $5,
			updated_at = now()
		WHERE id = $1
		  AND deleted_at IS NULL
		RETURNING
			id,
			shop_id,
			owner_id,
			name,
			category,
			price,
			detail,
			created_at,
			updated_at
	`,
		req.ID,
		req.Name,
		req.Category,
		req.Price,
		req.Detail,
	)

	var p domain.Product
	if err := row.Scan(
		&p.ID,
		&p.ShopID,
		&p.OwnerID,
		&p.Name,
		&p.Category,
		&p.Price,
		&p.Detail,
		&p.CreatedAt,
		&p.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			r.logger.WarnContext(ctx, "product not found for update",
				"productID", req.ID,
			)
			return nil, sharedErr.ErrNotFound
		}
		r.logger.ErrorContext(ctx, "failed to update product",
			"error", err,
			"productID", req.ID,
		)
		return nil, err
	}

	r.logger.InfoContext(ctx, "product updated successfully",
		"productID", p.ID,
		"shopID", p.ShopID,
	)

	return &p, nil
}

func (r *PostgresProductRepository) Delete(
	ctx context.Context,
	id string,
) error {

	r.logger.DebugContext(ctx, "deleting product",
		"productID", id,
	)

	res, err := r.db.ExecContext(ctx, `
		UPDATE products
		SET deleted_at = now(),
		    updated_at = now()
		WHERE id = $1
		  AND deleted_at IS NULL
	`, id)

	if err != nil {
		r.logger.ErrorContext(ctx, "failed to delete product",
			"error", err,
			"productID", id,
		)
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		r.logger.WarnContext(ctx, "product not found for deletion",
			"productID", id,
		)
		return sharedErr.ErrNotFound
	}

	r.logger.InfoContext(ctx, "product deleted successfully",
		"productID", id,
	)

	return nil
}
