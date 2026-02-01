package handler

import (
	"context"

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
	// Call service to validate token
	resp, err := h.svc.ValidateToken(ctx, req)
	if err != nil {
		return nil, err
	}

	// Validate response
	if resp == nil {
		return nil, err
	}

	return resp, nil
}
