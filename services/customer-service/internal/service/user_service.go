package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	domain "customerservice/internal/domain/entities"
	"customerservice/proto/userpb"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserService struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{db: db}
}

func MustGetUserID(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "missing metadata")
	}

	userIDs := md.Get("x-user-id")
	if len(userIDs) == 0 || userIDs[0] == "" {
		return "", status.Error(codes.Unauthenticated, "user not authenticated")
	}

	return userIDs[0], nil
}

// GetUser retrieves user details by user ID
func (s *UserService) GetUserDetail(ctx context.Context) (*userpb.UserDetailResponse, error) {
	// md, ok := metadata.FromIncomingContext(ctx)
	// if !ok {
	// 	return nil, fmt.Errorf("no metadata in context")
	// }
	// fmt.Printf("md %v", userID)

	userID, err := MustGetUserID(ctx)
	// userIDs := md.Get("x-user-id")

	// if req.UserId == "" {
	// 	return nil, errors.New("user_id is required")
	// }

	// userID, err := uuid.Parse("")
	fmt.Printf("userID %v", userID)
	if err != nil {
		return nil, errors.New("invalid user_id format")
	}

	query := `
		SELECT
			id, name, username, email, bio,
			twofa_enabled, is_verified, email_verified_at,
			status, last_login, role_id, created_at, updated_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`

	user := domain.User{}
	err = s.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.Name,
		&user.Username,
		&user.Email,
		&user.Bio,
		&user.TwoFAEnabled,
		&user.IsVerified,
		&user.EmailVerifiedAt,
		&user.Status,
		&user.LastLogin,
		&user.RoleID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	// Convert to proto response
	response := &userpb.UserDetailResponse{
		Id:         user.ID.String(),
		Name:       user.Name,
		Username:   user.Username,
		Email:      user.Email,
		Bio:        stringPtr(user.Bio),
		IsVerified: user.IsVerified,
		Status:     user.Status,
		CreatedAt:  timestamppb.New(user.CreatedAt),
		UpdatedAt:  timestamppb.New(user.UpdatedAt),
	}

	if user.EmailVerifiedAt != nil {
		response.EmailVerifiedAt = timestamppb.New(*user.EmailVerifiedAt)
	}
	if user.LastLogin != nil {
		response.LastLogin = timestamppb.New(*user.LastLogin)
	}
	if user.RoleID != nil {
		response.RoleId = user.RoleID.String()
	}

	return response, nil
}

// UpdateUsername updates user's username
func (s *UserService) UpdateUsername(ctx context.Context, req *userpb.UpdateUsernameRequest) (*userpb.UpdateUsernameResponse, error) {
	if req.UserId == "" {
		return nil, errors.New("user_id is required")
	}
	if req.NewUsername == "" {
		return nil, errors.New("new_username is required")
	}
	if len(req.NewUsername) < 3 || len(req.NewUsername) > 255 {
		return nil, errors.New("username must be between 3 and 255 characters")
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, errors.New("invalid user_id format")
	}

	// Check if username already exists
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 AND id != $2 AND deleted_at IS NULL)`
	err = s.db.QueryRowContext(ctx, checkQuery, req.NewUsername, userID).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("username already exists")
	}

	// Update username
	now := time.Now()
	updateQuery := `
		UPDATE users
		SET username = $1, updated_at = $2
		WHERE id = $3 AND deleted_at IS NULL
		RETURNING updated_at
	`

	var updatedAt time.Time
	err = s.db.QueryRowContext(ctx, updateQuery, req.NewUsername, now, userID).Scan(&updatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	response := &userpb.UpdateUsernameResponse{
		Success:     true,
		Message:     "Username updated successfully",
		UserId:      req.UserId,
		NewUsername: req.NewUsername,
		UpdatedAt:   timestamppb.New(updatedAt),
	}

	return response, nil
}

// Helper function to convert *string to string pointer
func stringPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
