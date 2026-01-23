package domain

import (
	"time"

	"github.com/google/uuid"
)

type Shop struct {
	ID      uuid.UUID `db:"id" json:"id"`
	OwnerID uuid.UUID `db:"owner_id" json:"owner_id"`

	Name        string  `db:"name" json:"name"`
	Slug        string  `db:"slug" json:"slug"`
	Description *string `db:"description" json:"description,omitempty"`
	Logo        *string `db:"logo" json:"logo,omitempty"`

	IsActive bool `db:"is_active" json:"is_active"`

	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}
