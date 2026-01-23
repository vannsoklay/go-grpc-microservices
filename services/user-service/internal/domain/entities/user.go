package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                 uuid.UUID  `db:"id" json:"id"`
	Name               string     `db:"name" json:"name"`
	Username           string     `db:"username" json:"username"`
	Email              string     `db:"email" json:"email"`
	PasswordHash       string     `db:"password_hash" json:"-"` // Never expose
	Bio                *string    `db:"bio" json:"bio,omitempty"`
	TwoFASecret        *string    `db:"twofa_secret" json:"-"` // Never expose
	TwoFAEnabled       bool       `db:"twofa_enabled" json:"twofa_enabled"`
	IsVerified         bool       `db:"is_verified" json:"is_verified"`
	EmailVerifiedAt    *time.Time `db:"email_verified_at" json:"email_verified_at,omitempty"`
	Status             string     `db:"status" json:"status"` // active, inactive, suspended
	LastLogin          *time.Time `db:"last_login" json:"last_login,omitempty"`
	LastPasswordChange *time.Time `db:"last_password_change" json:"last_password_change,omitempty"`
	RoleID             *uuid.UUID `db:"role_id" json:"role_id,omitempty"`
	CreatedAt          time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt          *time.Time `db:"deleted_at" json:"deleted_at,omitempty"` // Soft delete
}
