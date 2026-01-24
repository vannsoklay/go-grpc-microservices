package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	err "hpkg/constants/responses"
	auth "hpkg/grpc/middeware"

	"authservice/internal/domain"
	"authservice/proto/authpb"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	db         *sql.DB
	jwtService *auth.JWTService
}

func NewAuthService(db *sql.DB, jwtService *auth.JWTService) *AuthService {
	return &AuthService{db, jwtService}
}

func (s *AuthService) Register(ctx context.Context, req *authpb.RegisterReq) (*authpb.RegsiterResp, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Insert user into database
	var userID string
	err = s.db.QueryRowContext(ctx,
		`INSERT INTO users (name, username, email, password_hash)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id`,
		req.Name, req.Username, req.Email, hashedPassword,
	).Scan(&userID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("failed to create user")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Generate tokens
	accessToken, err := s.jwtService.GenerateAccessToken(userID, "user")
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Build response
	user := &authpb.User{
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
	}

	resp := &authpb.RegsiterResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}

	return resp, nil
}

func (s *AuthService) Login(
	ctx context.Context,
	email string,
	password string,
) (*authpb.LoginResp, error) {

	user, err := s.FindUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid credentials")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Compare password hash
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate tokens
	accessToken, err := s.jwtService.GenerateAccessToken(user.ID, "shop_owner")
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &authpb.LoginResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: &authpb.User{
			Name:     user.Name,
			Username: user.Username,
			Email:    user.Email,
		},
	}, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, req *authpb.TokenReq) (*authpb.ValidateTokenResp, error) {
	if req == nil || req.Token == "" {
		return nil, err.ValidationServiceError(err.ErrTokenInvalidMsg)
	}

	claim, err := s.jwtService.ValidateAccessToken(req.Token)
	if err != nil || claim == nil {
		return nil, err
	}

	role, perms, err := s.FindRoleAndPerms(claim.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return &authpb.ValidateTokenResp{
		UserId:      claim.UserID,
		Role:        role,
		Permissions: perms,
	}, nil
}

func (s *AuthService) FindUserByEmail(
	ctx context.Context,
	email string,
) (*domain.UserRsp, error) {

	query := `
		SELECT
			id,
			name,
			username,
			email,
			password_hash
		FROM users
		WHERE email = $1
		  AND deleted_at IS NULL
	`

	var user domain.UserRsp

	err := s.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Username,
		&user.Email,
		&user.Password,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("query user by email failed: %w", err)
	}

	fmt.Printf("user: %+v\n", user)

	return &user, nil
}

func (s *AuthService) FindRoleAndPerms(
	userId string,
) (role string, perms []string, err error) {

	query := `
		SELECT
			r.name AS role_name,
			p.name AS permission_name
		FROM users u
		JOIN roles r ON u.role_id = r.id
		JOIN role_permissions rp ON rp.role_id = r.id
		JOIN permissions p ON p.id = rp.permission_id
		WHERE u.id = $1
	`

	rows, err := s.db.QueryContext(context.Background(), query, userId)
	fmt.Printf("data %v", rows)
	if err != nil {
		return "", nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var perm string

		if err = rows.Scan(&role, &perm); err != nil {
			return "", nil, err
		}

		perms = append(perms, perm)
	}

	if err = rows.Err(); err != nil {
		return "", nil, err
	}

	if role == "" {
		return "", nil, sql.ErrNoRows
	}

	return role, perms, nil
}
