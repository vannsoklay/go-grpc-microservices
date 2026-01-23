package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"userservice/internal/domain/dto"
)

type PostgreUserRepository struct {
	db *sql.DB
}

func (r *PostgreUserRepository) User(ctx context.Context, userID string) (*dto.UserDetailDTO, error) {
	var user dto.UserDetailDTO
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, username, email FROM users WHERE id = $1`,
		userID,
	).Scan(&user.ID, &user.Name, &user.Username, &user.Email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// func (r *PostgreAuthRepository) Register(ctx context.Context, req dto.RegisterReq) (*dto.RegisterRsp, error) {
// 	// Hash password
// 	hashedPassword, err := bcrypt.GenerateFromPassword(
// 		[]byte(req.Password),
// 		bcrypt.DefaultCost,
// 	)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to hash password: %w", err)
// 	}

// 	// Insert user into database
// 	var userID string
// 	err = r.db.QueryRowContext(ctx,
// 		`INSERT INTO users (name, username, email, password_hash)
// 		 VALUES ($1, $2, $3, $4)
// 		 RETURNING id`,
// 		req.Name, req.Username, req.Email, hashedPassword,
// 	).Scan(&userID)

// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, errors.New("failed to create user")
// 		}
// 		return nil, fmt.Errorf("database error: %w", err)
// 	}

// 	// Generate tokens
// 	accessToken, err := r.jwtService.GenerateAccessToken(userID, "user")
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to generate access token: %w", err)
// 	}

// 	refreshToken, err := r.jwtService.GenerateRefreshToken(userID)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
// 	}

// 	// Build response
// 	user := domain.UserRsp{
// 		ID:       userID,
// 		Name:     req.Name,
// 		Username: req.Username,
// 		Email:    req.Email,
// 	}

// 	resp := &dto.RegisterRsp{
// 		AccessToken:  accessToken,
// 		RefreshToken: refreshToken,
// 		User:         user,
// 	}

// 	return resp, nil
// }
