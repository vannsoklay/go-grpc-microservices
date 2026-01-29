package domain

import "time"

type Product struct {
	ID          string     `db:"id"`
	ShopID      string     `db:"shop_id"`
	OwnerID     string     `db:"owner_id"`
	Name        string     `db:"name"`
	Category    string     `db:"category"`
	Price       float64    `db:"price"`
	Description string     `db:"description"`
	Detail      *string    `db:"detail"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
}

type CreateProductRequest struct {
	ShopID      string
	OwnerID     string
	Name        string
	Description string
	Category    string
	Price       float64
	Detail      string
}

type UpdateProductRequest struct {
	ID          string
	Name        string
	Category    string
	Price       float64
	Description string
	Detail      string
}
