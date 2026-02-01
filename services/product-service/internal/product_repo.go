package repository

import (
	"context"
	"productservice/internal/domain"
)

type ProductRepository interface {
	Create(ctx context.Context, params domain.CreateProductRequest) (*domain.Product, error)
	GetByID(ctx context.Context, id string) (*domain.Product, error)
	Update(ctx context.Context, params domain.UpdateProductRequest) (*domain.Product, error)
	Delete(ctx context.Context, id string) error
}
