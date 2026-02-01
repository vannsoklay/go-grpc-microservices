package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	pagination "hpkg/constants"
	"productservice/internal/domain"
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
) ([]*domain.Product, string, int64, int64, error) {

	if limit <= 0 || limit > 50 {
		limit = 20
	}

	// -----------------------------------
	// Base WHERE (shared)
	// -----------------------------------
	baseWhere := `
		FROM products
		WHERE shop_id = $1
		  AND deleted_at IS NULL
	`

	args := []any{shopID}
	argPos := 2

	// Search
	if search != "" {
		baseWhere += fmt.Sprintf(" AND name ILIKE $%d", argPos)
		args = append(args, "%"+search+"%")
		argPos++
	}

	// Filter
	if filter != "" {
		baseWhere += fmt.Sprintf(" AND LOWER(category) = LOWER($%d)", argPos)
		args = append(args, filter)
		argPos++
	}

	// -----------------------------------
	// TOTAL ALL PRODUCTS (no filter)
	// -----------------------------------
	var totalAllCount int64
	if err := r.db.QueryRowContext(
		ctx,
		`
		SELECT COUNT(1)
		FROM products
		WHERE shop_id = $1
		  AND deleted_at IS NULL
		`,
		shopID,
	).Scan(&totalAllCount); err != nil {
		r.logger.ErrorContext(ctx, "failed to count all products",
			"error", err,
			"shopID", shopID,
		)
		return nil, "", 0, 0, err
	}

	// -----------------------------------
	// TOTAL FILTERED PRODUCTS
	// -----------------------------------
	var totalCount int64
	countQuery := `SELECT COUNT(1) ` + baseWhere

	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount); err != nil {
		r.logger.ErrorContext(ctx, "failed to count filtered products",
			"error", err,
			"shopID", shopID,
		)
		return nil, "", 0, 0, err
	}

	// -----------------------------------
	// LIST QUERY (cursor pagination)
	// -----------------------------------
	listArgs := append([]any{}, args...)
	listArgPos := argPos

	listQuery := `
		SELECT id, shop_id, owner_id, name, category, price,
		       description, detail, created_at, updated_at, deleted_at
	` + baseWhere

	if cursor != "" {
		c, err := pagination.DecodeCursor(cursor)
		if err != nil {
			r.logger.ErrorContext(ctx, "failed to decode cursor",
				"error", err,
				"cursor", cursor,
			)
			return nil, "", 0, 0, err
		}

		// forward-only cursor
		listQuery += fmt.Sprintf(
			" AND (created_at, id) < ($%d, $%d)",
			listArgPos,
			listArgPos+1,
		)
		listArgs = append(listArgs, c.CreatedAt, c.ID)
		listArgPos += 2
	}

	order := "ASC"
	if sortDesc {
		order = "DESC"
	}

	listQuery += fmt.Sprintf(
		" ORDER BY %s %s, id %s LIMIT $%d",
		sortColumn,
		order,
		order,
		listArgPos,
	)
	listArgs = append(listArgs, limit+1)

	rows, err := r.db.QueryContext(ctx, listQuery, listArgs...)
	if err != nil {
		r.logger.ErrorContext(ctx, "failed to query products",
			"error", err,
			"query", listQuery,
			"args", listArgs,
		)
		return nil, "", 0, 0, err
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(
			&p.ID, &p.ShopID, &p.OwnerID,
			&p.Name, &p.Category, &p.Price,
			&p.Description, &p.Detail,
			&p.CreatedAt, &p.UpdatedAt, &p.DeletedAt,
		); err != nil {
			return nil, "", 0, 0, err
		}
		products = append(products, &p)
	}

	if err := rows.Err(); err != nil {
		return nil, "", 0, 0, err
	}

	// -----------------------------------
	// NEXT CURSOR
	// -----------------------------------
	var nextCursor string
	if len(products) > limit {
		last := products[limit-1]
		c, err := pagination.EncodeCursor(pagination.ProductCursor{
			CreatedAt: last.CreatedAt,
			ID:        last.ID,
		})
		if err != nil {
			return nil, "", 0, 0, err
		}
		nextCursor = c
		products = products[:limit]
	}

	return products, nextCursor, totalCount, totalAllCount, nil
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
		&p.Description,
		&p.Detail,
		&p.CreatedAt,
		&p.UpdatedAt,
		&p.DeletedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			r.logger.WarnContext(ctx, "product not found",
				"productID", id,
			)
			return nil, err
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
		"price", req.Price,
		"category", req.Description,
		"detail", req.Detail,
	)

	row := r.db.QueryRowContext(ctx, `
		INSERT INTO products (
			shop_id,
			owner_id,
			name,
			price,
			category,
			description,
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
			description,
			detail,
			created_at,
			updated_at
	`,
		req.ShopID,
		req.OwnerID,
		req.Name,
		req.Price,
		req.Category,
		req.Description,
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
		&p.Description,
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
			description = $5,
			detail = $6,
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
		&p.Description,
		&p.Detail,
		&p.CreatedAt,
		&p.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			r.logger.WarnContext(ctx, "product not found for update",
				"productID", req.ID,
			)
			return nil, err
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

	rows, err := res.RowsAffected()
	if rows == 0 {
		r.logger.WarnContext(ctx, "product not found for deletion",
			"productID", id,
		)
		return err
	}

	r.logger.InfoContext(ctx, "product deleted successfully",
		"productID", id,
	)

	return nil
}
