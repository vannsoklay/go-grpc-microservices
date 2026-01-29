package service

import (
	"context"
	"database/sql"
	"errors"

	errs "hpkg/constants/responses"
	pkg "hpkg/grpc"
	proto "userservice/internal/domain/proto"
	"userservice/internal/repository"
	"userservice/proto/userpb"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// ---------------------------
// Get user details
// ---------------------------
func (s *UserService) GetUserDetail(ctx context.Context) (*userpb.UserDetailResponse, error) {
	userID, err := pkg.MustGetUserID(ctx)
	if err != nil {
		return nil, errs.GRPC(codes.Unauthenticated, errs.UnauthenticatedCode, errs.UnauthenticatedMsg)
	}

	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, errs.GRPC(codes.Canceled, errs.RequestCanceledCode, errs.RequestCanceledMsg)
		}
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.GRPC(codes.NotFound, errs.UserNotFoundCode, errs.UserNotFoundMsg)
		}
		return nil, errs.GRPC(codes.Internal, errs.UserFetchFailedCode, errs.UserFetchFailedMsg)
	}

	return proto.MapUserToProto(user), nil
}

// ---------------------------
// Update username
// ---------------------------
func (s *UserService) UpdateUsername(ctx context.Context, req *userpb.UpdateUsernameRequest) (*userpb.UpdateUsernameResponse, error) {
	if req.UserId == "" {
		return nil, errs.GRPC(codes.InvalidArgument, errs.InvalidRequestCode, errs.InvalidRequestMsg)
	}
	if req.NewUsername == "" {
		return nil, errs.GRPC(codes.InvalidArgument, errs.InvalidRequestCode, errs.InvalidRequestMsg)
	}
	if len(req.NewUsername) < 3 || len(req.NewUsername) > 255 {
		return nil, errs.GRPC(codes.InvalidArgument, errs.InvalidRequestCode, errs.InvalidRequestMsg)
	}

	userUUID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, errs.GRPC(codes.InvalidArgument, errs.InvalidRequestCode, errs.InvalidRequestMsg)
	}

	// Check if username exists
	exists, err := s.repo.IsUsernameExists(ctx, req.NewUsername, userUUID.String())
	if err != nil {
		return nil, errs.GRPC(codes.Internal, errs.UsernameCheckFailedCode, errs.UsernameCheckFailedMsg)
	}
	if exists {
		return nil, errs.GRPC(codes.AlreadyExists, errs.UsernameExistsCode, errs.UsernameExistsMsg)
	}

	updatedAt, err := s.repo.UpdateUsername(ctx, userUUID.String(), req.NewUsername)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.GRPC(codes.NotFound, errs.UserNotFoundCode, errs.UserNotFoundMsg)
		}
		return nil, errs.GRPC(codes.Internal, errs.UsernameUpdateFailedCode, errs.UsernameUpdateFailedMsg)
	}

	return &userpb.UpdateUsernameResponse{
		Success:     true,
		Message:     "Username updated successfully",
		UserId:      req.UserId,
		NewUsername: req.NewUsername,
		UpdatedAt:   timestamppb.New(updatedAt),
	}, nil
}
