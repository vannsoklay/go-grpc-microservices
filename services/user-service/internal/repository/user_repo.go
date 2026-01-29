package repository

import (
	"context"
	"time"
	domain "userservice/internal/domain/entities"
)

type UserRepository interface {
	GetByID(ctx context.Context, id string) (*domain.User, error)
	UpdateUsername(ctx context.Context, id string, newUsername string) (time.Time, error)
	IsUsernameExists(ctx context.Context, username, excludeID string) (bool, error)
}
