package handler

import (
	"context"
	"userservice/internal/service"
	"userservice/proto/userpb"

	"google.golang.org/protobuf/types/known/emptypb"
)

type UserHandler struct {
	userpb.UnimplementedUserServiceServer
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) GetUserDetail(
	ctx context.Context,
	_ *emptypb.Empty,
) (*userpb.UserDetailResponse, error) {
	resp, err := h.svc.GetUserDetail(ctx)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *UserHandler) UpdateUsername(
	ctx context.Context,
	req *userpb.UpdateUsernameRequest,
) (*userpb.UpdateUsernameResponse, error) {
	resp, err := h.svc.UpdateUsername(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
