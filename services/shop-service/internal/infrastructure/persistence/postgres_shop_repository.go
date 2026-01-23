package persistence

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"shopservice/internal/domain/dto"
)

type ShopRepository interface {
	CreateShop(ctx context.Context, shop *dto.ShopDTO) (string, error)
	GetByOwnerID(ctx context.Context, ownerID string) (*dto.ShopDTO, error)
	GetBySlug(ctx context.Context, slug string) (bool, error)
	UpdateShop(ctx context.Context, shop *dto.ShopDTO) (*dto.ShopDTO, error)
	DeleteByOwnerID(ctx context.Context, ownerID string) (int64, error)
}

type PostgresShopRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewPostgresShopRepository(db *sql.DB, logger *slog.Logger) *PostgresShopRepository {
	return &PostgresShopRepository{
		db:     db,
		logger: logger,
	}
}

const (
	queryShopByOwnerID = `
		SELECT id, owner_id, name, slug, description, logo, is_active, created_at, updated_at
		FROM shops
		WHERE owner_id = $1 AND deleted_at IS NULL
	`
	querySlugExists = `
		SELECT EXISTS (SELECT 1 FROM shops WHERE slug = $1 AND deleted_at IS NULL)
	`
	queryCreateShop = `
		INSERT INTO shops (id, owner_id, name, slug, description, logo, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, true, $7, $7)
	`
	queryUpdateShop = `
		UPDATE shops
		SET name = $1, description = $2, logo = $3, is_active = $4, updated_at = now()
		WHERE owner_id = $5 AND deleted_at IS NULL
		RETURNING id, owner_id, name, slug, description, logo, is_active, created_at, updated_at
	`
	queryDeleteShop = `
		UPDATE shops SET deleted_at = now()
		WHERE owner_id = $1 AND deleted_at IS NULL
	`
)

func (r *PostgresShopRepository) CreateShop(ctx context.Context, shop *dto.ShopDTO) (string, error) {
	_, err := r.db.ExecContext(ctx, queryCreateShop,
		shop.ID, shop.OwnerID, shop.Name, shop.Slug,
		nullStr(shop.Description), nullStr(shop.Logo), shop.CreatedAt,
	)
	if err != nil {
		r.logger.ErrorContext(ctx, "failed to create shop",
			slog.String("owner_id", shop.OwnerID),
			slog.String("slug", shop.Slug),
			slog.String("error", err.Error()),
		)
		return "", err
	}
	r.logger.InfoContext(ctx, "shop created successfully",
		slog.String("shop_id", shop.ID),
		slog.String("owner_id", shop.OwnerID),
	)
	return shop.ID, nil
}

func (r *PostgresShopRepository) GetByOwnerID(ctx context.Context, ownerID string) (*dto.ShopDTO, error) {
	row := r.db.QueryRowContext(ctx, queryShopByOwnerID, ownerID)
	shop, err := scanShop(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.DebugContext(ctx, "shop not found",
				slog.String("owner_id", ownerID),
			)
			return nil, err
		}
		r.logger.ErrorContext(ctx, "failed to query shop",
			slog.String("owner_id", ownerID),
			slog.String("error", err.Error()),
		)
		return nil, err
	}
	return shop, nil
}

func (r *PostgresShopRepository) GetBySlug(ctx context.Context, slug string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, querySlugExists, slug).Scan(&exists)
	if err != nil {
		r.logger.ErrorContext(ctx, "failed to check slug uniqueness",
			slog.String("slug", slug),
			slog.String("error", err.Error()),
		)
		return false, err
	}
	return exists, nil
}

func (r *PostgresShopRepository) UpdateShop(ctx context.Context, shop *dto.ShopDTO) (*dto.ShopDTO, error) {
	row := r.db.QueryRowContext(ctx, queryUpdateShop,
		shop.Name, nullStr(shop.Description), nullStr(shop.Logo), shop.IsActive, shop.OwnerID,
	)
	updated, err := scanShop(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.WarnContext(ctx, "shop not found for update",
				slog.String("owner_id", shop.OwnerID),
			)
			return nil, err
		}
		r.logger.ErrorContext(ctx, "failed to update shop",
			slog.String("owner_id", shop.OwnerID),
			slog.String("error", err.Error()),
		)
		return nil, err
	}
	r.logger.InfoContext(ctx, "shop updated successfully",
		slog.String("shop_id", updated.ID),
	)
	return updated, nil
}

func (r *PostgresShopRepository) DeleteByOwnerID(ctx context.Context, ownerID string) (int64, error) {
	res, err := r.db.ExecContext(ctx, queryDeleteShop, ownerID)
	if err != nil {
		r.logger.ErrorContext(ctx, "failed to delete shop",
			slog.String("owner_id", ownerID),
			slog.String("error", err.Error()),
		)
		return 0, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		r.logger.ErrorContext(ctx, "failed to get rows affected",
			slog.String("owner_id", ownerID),
			slog.String("error", err.Error()),
		)
		return 0, err
	}
	if affected > 0 {
		r.logger.InfoContext(ctx, "shop deleted successfully",
			slog.String("owner_id", ownerID),
		)
	}
	return affected, nil
}

func scanShop(row interface{ Scan(...interface{}) error }) (*dto.ShopDTO, error) {
	var s dto.ShopDTO
	err := row.Scan(
		&s.ID, &s.OwnerID, &s.Name, &s.Slug,
		&s.Description, &s.Logo, &s.IsActive, &s.CreatedAt, &s.UpdatedAt,
	)
	return &s, err
}

func nullStr(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *s, Valid: true}
}
