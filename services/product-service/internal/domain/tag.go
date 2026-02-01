package domain

import (
	"time"
)

type Tag struct {
	ID        string     `db:"id"`
	ShopID    string     `db:"shop_id"`
	Name      string     `db:"name"`
	Slug      string     `db:"slug"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

// TagWithCount is used for the "List" view (Tag + metadata)
type TagWithCount struct {
	Tag
	ProductCount int32 `db:"product_count"`
}

// TagDetail is used for the "Detail" view (Tag + associated products)
type TagDetail struct {
	Tag
	Products      []*ProductSummary
	TotalProducts int32
}
