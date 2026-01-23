package repository

import (
	"context"
	"paymentservice/internal/domain"
)

type PaymentRepository interface {
	Create(ctx context.Context, p *domain.Payment) (*domain.Payment, error)
	GetByID(ctx context.Context, id string) (*domain.Payment, error)
	ListByUser(ctx context.Context, userID string, page, size string, status string) ([]*domain.Payment, int32, error)
	ListByOrder(ctx context.Context, orderID string, page, size string) ([]*domain.Payment, int32, error)
	UpdateStatus(ctx context.Context, id string, status string) error
}
