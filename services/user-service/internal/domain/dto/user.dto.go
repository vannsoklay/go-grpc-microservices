package dto

import (
	"time"

	"github.com/google/uuid"
)

type UserDetailDTO struct {
	ID              uuid.UUID  `json:"id"`
	Name            string     `json:"name"`
	Username        string     `json:"username"`
	Email           string     `json:"email"`
	Bio             *string    `json:"bio,omitempty"`
	TwoFAEnabled    bool       `json:"twofa_enabled"`
	IsVerified      bool       `json:"is_verified"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
	Status          string     `json:"status"`
	LastLogin       *time.Time `json:"last_login,omitempty"`
	RoleID          *uuid.UUID `json:"role_id,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
