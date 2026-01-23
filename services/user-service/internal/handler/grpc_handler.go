package handler

import (
	"context"
	"fmt"
	"log"
	"userservice/internal/service"
	"userservice/proto/userpb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	fmt.Printf("service it works")
	resp, err := h.svc.GetUserDetail(ctx)
	if err != nil {
		log.Printf("Error getting user detail: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return resp, nil
}

func (h *UserHandler) UpdateUsername(
	ctx context.Context,
	req *userpb.UpdateUsernameRequest,
) (*userpb.UpdateUsernameResponse, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if req.NewUsername == "" {
		return nil, status.Error(codes.InvalidArgument, "new_username is required")
	}

	resp, err := h.svc.UpdateUsername(ctx, req)
	if err != nil {
		log.Printf("Error updating username: %v", err)

		// Map service errors to gRPC status codes
		switch err.Error() {
		case "user not found":
			return nil, status.Error(codes.NotFound, "user not found")
		case "username already exists":
			return nil, status.Error(codes.AlreadyExists, "username already exists")
		case "invalid user_id format":
			return nil, status.Error(codes.InvalidArgument, "invalid user_id format")
		default:
			return nil, status.Error(codes.Internal, "internal server error")
		}
	}

	return resp, nil
}
