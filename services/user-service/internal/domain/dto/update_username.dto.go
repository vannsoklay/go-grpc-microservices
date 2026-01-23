package dto

import (
	"time"

	"github.com/google/uuid"
)

type UpdateUsernameRequest struct {
	UserID      uuid.UUID `json:"user_id" binding:"required"`
	NewUsername string    `json:"new_username" binding:"required,min=3,max=255"`
}

type UpdateUsernameResponse struct {
	Success     bool      `json:"success"`
	Message     string    `json:"message"`
	UserID      uuid.UUID `json:"user_id"`
	NewUsername string    `json:"new_username"`
	UpdatedAt   time.Time `json:"updated_at"`
}
