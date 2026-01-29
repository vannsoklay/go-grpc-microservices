package repository

import (
	"context"
	"database/sql"
	"time"

	"log/slog"
	domain "userservice/internal/domain/entities"
)

type PostgresUserRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewPostgresUserRepository(db *sql.DB, logger *slog.Logger) *PostgresUserRepository {
	return &PostgresUserRepository{
		db:     db,
		logger: logger,
	}
}

// ---------------------------
// GET USER BY ID
// ---------------------------
func (r *PostgresUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	r.logger.DebugContext(ctx, "fetching user by ID", "userID", id)

	query := `
		SELECT id, name, username, email, bio, twofa_enabled, is_verified,
		       email_verified_at, status, last_login, role_id, created_at, updated_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Name, &user.Username, &user.Email, &user.Bio,
		&user.TwoFAEnabled, &user.IsVerified, &user.EmailVerifiedAt,
		&user.Status, &user.LastLogin, &user.RoleID, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.WarnContext(ctx, "user not found", "userID", id)
		} else {
			r.logger.ErrorContext(ctx, "failed to fetch user", "error", err, "userID", id)
		}
		return nil, err
	}

	r.logger.DebugContext(ctx, "user fetched successfully", "userID", id)
	return &user, nil
}

// ---------------------------
// UPDATE USERNAME
// ---------------------------
func (r *PostgresUserRepository) UpdateUsername(ctx context.Context, id string, newUsername string) (time.Time, error) {
	now := time.Now()
	r.logger.DebugContext(ctx, "updating username", "userID", id, "newUsername", newUsername)

	query := `
		UPDATE users
		SET username = $1, updated_at = $2
		WHERE id = $3 AND deleted_at IS NULL
		RETURNING updated_at
	`
	var updatedAt time.Time
	err := r.db.QueryRowContext(ctx, query, newUsername, now, id).Scan(&updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.WarnContext(ctx, "user not found for update", "userID", id)
		} else {
			r.logger.ErrorContext(ctx, "failed to update username", "error", err, "userID", id)
		}
		return time.Time{}, err
	}

	r.logger.InfoContext(ctx, "username updated successfully", "userID", id, "newUsername", newUsername, "updatedAt", updatedAt)
	return updatedAt, nil
}

// ---------------------------
// CHECK IF USERNAME EXISTS
// ---------------------------
func (r *PostgresUserRepository) IsUsernameExists(ctx context.Context, username, excludeID string) (bool, error) {
	r.logger.DebugContext(ctx, "checking if username exists", "username", username, "excludeID", excludeID)

	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 AND id != $2 AND deleted_at IS NULL)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, username, excludeID).Scan(&exists)
	if err != nil {
		r.logger.ErrorContext(ctx, "failed to check username existence", "error", err, "username", username)
		return false, err
	}

	r.logger.DebugContext(ctx, "username existence check completed", "username", username, "exists", exists)
	return exists, nil
}
