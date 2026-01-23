package domain

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID     uuid.UUID `db:"id" json:"id"`
	ShopID uuid.UUID `db:"shop_id" json:"shop_id"`

	// Basic info
	Name  string  `db:"name" json:"name"`
	Email *string `db:"email" json:"email,omitempty"`
	Phone *string `db:"phone" json:"phone,omitempty"`

	// Optional profile
	Avatar *string `db:"avatar" json:"avatar,omitempty"`
	Note   *string `db:"note" json:"note,omitempty"`

	// Status
	IsActive bool `db:"is_active" json:"is_active"`

	// Timestamps
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}
