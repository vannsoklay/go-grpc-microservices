package service

import (
	"context"
	"database/sql"
	"fmt"
	Err "hpkg/constants/responses"
	auth "hpkg/grpc/middeware"

	"authservice/internal/domain"
	"authservice/proto/authpb"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
)

type AuthService struct {
	db         *sql.DB
	jwtService *auth.JWTService
}

func NewAuthService(db *sql.DB, jwtService *auth.JWTService) *AuthService {
	return &AuthService{db, jwtService}
}

func (s *AuthService) Register(
	ctx context.Context,
	req *authpb.RegisterReq,
) (*authpb.RegsiterResp, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, Err.GRPC(
			codes.Internal,
			Err.ErrUserCreateFailedCode,
			Err.ErrUserCreateFailedMsg,
		)
	}

	var userID string
	err = s.db.QueryRowContext(
		ctx,
		`INSERT INTO users (name, username, email, password_hash)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id`,
		req.Name, req.Username, req.Email, hashedPassword,
	).Scan(&userID)

	if err != nil {
		return nil, Err.GRPC(
			codes.Internal,
			Err.ErrUserCreateFailedCode,
			Err.ErrUserCreateFailedMsg,
		)
	}

	accessToken, err := s.jwtService.GenerateAccessToken(userID, "user")
	if err != nil {
		return nil, Err.GRPC(
			codes.Internal,
			Err.TokenGenerateFailedCode,
			Err.AccessTokenGenerateFailedMsg,
		)
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(userID)
	if err != nil {
		return nil, Err.GRPC(
			codes.Internal,
			Err.TokenGenerateFailedCode,
			Err.RefreshTokenGenerateFailedMsg,
		)
	}

	return &authpb.RegsiterResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: &authpb.User{
			Name:     req.Name,
			Username: req.Username,
			Email:    req.Email,
		},
	}, nil
}

func (s *AuthService) Login(
	ctx context.Context,
	email string,
	password string,
) (*authpb.LoginResp, error) {

	user, err := s.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, Err.GRPC(
			codes.Unauthenticated,
			Err.ErrInvalidCredentialsCode,
			Err.ErrInvalidCredentialsMsg,
		)
	}

	if bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	) != nil {
		return nil, Err.GRPC(
			codes.Unauthenticated,
			Err.ErrInvalidCredentialsCode,
			Err.ErrInvalidCredentialsMsg,
		)
	}

	accessToken, err := s.jwtService.GenerateAccessToken(user.ID, "shop_owner")
	if err != nil {
		return nil, Err.GRPC(
			codes.Internal,
			Err.TokenGenerateFailedCode,
			Err.AccessTokenGenerateFailedMsg,
		)
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, Err.GRPC(
			codes.Internal,
			Err.TokenGenerateFailedCode,
			Err.RefreshTokenGenerateFailedMsg,
		)
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
		return nil, Err.GRPC(codes.Unauthenticated, Err.ErrTokenInvalidCode, Err.ErrTokenInvalidMsg)
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
