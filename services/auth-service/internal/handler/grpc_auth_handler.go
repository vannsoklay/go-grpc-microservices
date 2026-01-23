package handler

import (
	"context"
	"fmt"
	"hpkg/errors"
	"log"

	"authservice/internal/service"
	"authservice/proto/authpb"
)

type AuthHandler struct {
	authpb.UnimplementedAuthServiceServer
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Register(
	ctx context.Context,
	req *authpb.RegisterReq,
) (*authpb.RegsiterResp, error) {
	resp, err := h.svc.Register(ctx, req)
	if err != nil {
		return &authpb.RegsiterResp{}, nil
	}

	return resp, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *authpb.LoginReq) (*authpb.LoginResp, error) {
	resp, err := h.svc.Login(ctx, *req.Email, req.Password)
	if err != nil {
		return &authpb.LoginResp{}, nil
	}

	return resp, nil
}

func (h *AuthHandler) Validate(ctx context.Context, req *authpb.TokenReq) (*authpb.ValidateTokenResp, error) {
	// Validate request
	fmt.Printf("data %v", req.Token)
	if err := h.validateRequest(req); err != nil {
		log.Printf("validation error: %v", err)
		return nil, err
	}

	// Call service to validate token
	resp, err := h.svc.ValidateToken(ctx, req)
	if err != nil {
		log.Printf("service error: %v", err)
		return nil, err
	}

	fmt.Printf("resp %v", resp)

	// Validate response
	if resp == nil {
		log.Println("service returned nil response")
		return nil, err
	}

	return resp, nil
}

// validateRequest validates the incoming request
func (h *AuthHandler) validateRequest(req *authpb.TokenReq) error {
	if req == nil {
		return errors.ErrInvalidInput
	}

	if req.Token == "" {
		return errors.ErrInvalidInput
	}

	return nil
}

// func (h *AuthHandler) GetPermissions(ctx context.Context, req *authpb.GetPermissionsRequest) (*authpb.GetPermissionsResponse, error) {
// 	resp, err := h.svc.GetPermissions(ctx, req.Role)
// 	if err != nil {
// 		return &authpb.GetPermissionsResponse{}, nil
// }

// func (h *AuthHandler) ValidateToken(
// 	ctx context.Context,
// 	req *authpb.ValidateTokenRequest,
// ) (*authpb.ValidateTokenResponse, error) {

// 	claims, err := h.svc.ValidateToken(req.Token)
// 	if err != nil {
// 		return &authpb.ValidateTokenResponse{
// 			Valid:   false,
// 			Message: err.Error(),
// 		}, nil
// 	}

// 	perms, _ := h.svc.GetPermissions(claims.Role)

// 	return &authpb.ValidateTokenResponse{
// 		Valid:       true,
// 		UserId:      claims.UserID,
// 		Username:    claims.Username,
// 		Role:        claims.Role,
// 		Permissions: perms,
// 		Message:     "valid token",
// 	}, nil
// }

// func (h *AuthHandler) GetPermissions(
// 	ctx context.Context,
// 	req *authpb.GetPermissionsRequest,
// ) (*authpb.GetPermissionsResponse, error) {

// 	perms, _ := h.svc.GetPermissions(req.Role)

// 	return &authpb.GetPermissionsResponse{
// 		Role:        req.Role,
// 		Permissions: perms,
// 	}, nil
// }

// func (h *AuthHandler) CheckPermission(
// 	ctx context.Context,
// 	req *authpb.CheckPermissionRequest,
// ) (*authpb.CheckPermissionResponse, error) {

// 	claims, err := h.svc.ValidateToken(req.Token)
// 	if err != nil {
// 		return &authpb.CheckPermissionResponse{HasPermission: false}, nil
// 	}

// 	ok := h.svc.CheckPermission(req.Role, req.Permission)

// 	return &authpb.CheckPermissionResponse{
// 		HasPermission: ok,
// 		UserId:        claims.UserID,
// 		Username:      claims.Username,
// 		Role:          claims.Role,
// 		Permission:    req.Permission,
// 	}, nil
// }

// func (h *AuthHandler) Logout(
// 	ctx context.Context,
// 	req *authpb.LogoutRequest,
// ) (*authpb.LogoutResponse, error) {

// 	h.svc.Logout(req.Token)

// 	return &authpb.LogoutResponse{
// 		Message:  "logout success",
// 		Username: req.Username,
// 	}, nil
// }

// func (h *AuthHandler) RefreshToken(
// 	ctx context.Context,
// 	req *authpb.RefreshTokenRequest,
// ) (*authpb.RefreshTokenResponse, error) {

// 	token, exp, err := h.svc.RefreshToken(req.Token)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &authpb.RefreshTokenResponse{
// 		Token:     token,
// 		ExpiresAt: exp,
// 	}, nil
// }
