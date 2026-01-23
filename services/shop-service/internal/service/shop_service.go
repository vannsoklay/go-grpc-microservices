package service

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"shopservice/internal/domain/dto"
	"shopservice/internal/infrastructure/persistence"
	"shopservice/proto/shoppb"
)

type ShopService struct {
	repo   persistence.ShopRepository
	logger *slog.Logger
}

func NewShopService(repo persistence.ShopRepository, logger *slog.Logger) *ShopService {
	return &ShopService{
		repo:   repo,
		logger: logger,
	}
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

func (s *ShopService) CreateShop(ctx context.Context, req *shoppb.CreateShopRequest) (*shoppb.CreateShopResponse, error) {
	if req.Name == "" || req.Slug == "" {
		return nil, status.Error(codes.InvalidArgument, "name and slug are required")
	}

	ownerID, err := MustGetUserID(ctx)
	if err != nil {
		return nil, err
	}

	// Check slug uniqueness
	exists, err := s.repo.GetBySlug(ctx, req.Slug)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to check slug uniqueness", slog.String("error", err.Error()))
		return nil, status.Error(codes.Internal, "failed to create shop")
	}
	if exists {
		s.logger.WarnContext(ctx, "slug already exists", slog.String("slug", req.Slug))
		return nil, status.Error(codes.AlreadyExists, "shop slug already exists")
	}

	now := time.Now()
	shop := &dto.ShopDTO{
		ID:          uuid.New().String(),
		OwnerID:     ownerID,
		Name:        req.Name,
		Slug:        req.Slug,
		Description: emptyStrToNil(req.Description),
		Logo:        emptyStrToNil(req.Logo),
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	id, err := s.repo.CreateShop(ctx, shop)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create shop")
	}

	return &shoppb.CreateShopResponse{
		ShopId:    id,
		CreatedAt: timestamppb.New(now),
	}, nil
}

func (s *ShopService) GetMyShop(ctx context.Context, _ *shoppb.GetMyShopRequest) (*shoppb.ShopResponse, error) {
	ownerID, err := MustGetUserID(ctx)
	if err != nil {
		return nil, err
	}

	shop, err := s.repo.GetByOwnerID(ctx, ownerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "shop not found")
		}
		return nil, status.Error(codes.Internal, "failed to retrieve shop")
	}

	return toShopResponse(shop), nil
}

func (s *ShopService) UpdateShop(ctx context.Context, req *shoppb.UpdateShopRequest) (*shoppb.ShopResponse, error) {
	ownerID, err := MustGetUserID(ctx)
	if err != nil {
		return nil, err
	}

	shop := &dto.ShopDTO{
		OwnerID:     ownerID,
		Name:        req.Name,
		Description: emptyStrToNil(req.Description),
		Logo:        emptyStrToNil(req.Logo),
		IsActive:    req.IsActive,
	}

	updated, err := s.repo.UpdateShop(ctx, shop)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "shop not found")
		}
		return nil, status.Error(codes.Internal, "failed to update shop")
	}

	return toShopResponse(updated), nil
}

func (s *ShopService) DeleteShop(ctx context.Context, _ *shoppb.DeleteShopRequest) (*shoppb.DeleteShopResponse, error) {
	ownerID, err := MustGetUserID(ctx)
	if err != nil {
		return nil, err
	}

	affected, err := s.repo.DeleteByOwnerID(ctx, ownerID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to delete shop")
	}

	if affected == 0 {
		return nil, status.Error(codes.NotFound, "shop not found")
	}

	return &shoppb.DeleteShopResponse{Success: true}, nil
}

func toShopResponse(s *dto.ShopDTO) *shoppb.ShopResponse {
	return &shoppb.ShopResponse{
		Id:          s.ID,
		OwnerId:     s.OwnerID,
		Name:        s.Name,
		Slug:        s.Slug,
		Description: ptrOrEmpty(s.Description),
		Logo:        ptrOrEmpty(s.Logo),
		IsActive:    s.IsActive,
		CreatedAt:   timestamppb.New(s.CreatedAt),
		UpdatedAt:   timestamppb.New(s.UpdatedAt),
	}
}

func emptyStrToNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func ptrOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
