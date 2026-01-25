package service

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"hpkg/constants/responses"
	reqCtx "hpkg/grpc"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
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

func (s *ShopService) ValidateShop(ctx context.Context, req *shoppb.ValidateShopRequest) (*shoppb.ValidateShopResponse, error) {
	ownerID, err := reqCtx.MustGetUserID(ctx)
	if err != nil {
		return nil, err
	}

	id, slug, err := s.repo.ValidateShop(ctx, ownerID, req.ShopId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "shop not found")
		}
		return nil, status.Error(codes.Internal, "failed to retrieve shop")
	}

	return &shoppb.ValidateShopResponse{
		Id:   id,
		Slug: slug,
	}, nil
}

func (s *ShopService) CreateShop(ctx context.Context, req *shoppb.CreateShopRequest) (*shoppb.CreateShopResponse, error) {
	// Validate required fields
	if req.Name == "" || req.Slug == "" {
		return nil, status.Error(codes.InvalidArgument, "name and slug are required")
	}

	// Get current user ID from context
	ownerID, err := reqCtx.MustGetUserID(ctx)
	if err != nil {
		return nil, err
	}

	// Check how many shops the user already has
	count, err := s.repo.CountShopsByOwner(ctx, ownerID)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to count user's shops",
			slog.String("owner_id", ownerID),
			slog.String("error", err.Error()),
		)
		return nil, status.Error(codes.Internal, "failed to create shop")
	}

	if count >= 2 {
		s.logger.WarnContext(ctx, "user reached maximum allowed shops",
			slog.String("owner_id", ownerID),
		)
		return nil, responses.NewShopLimitError(count)
	}

	// Check slug uniqueness
	exists, err := s.repo.GetBySlug(ctx, req.Slug)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to check slug uniqueness",
			slog.String("error", err.Error()),
		)
		return nil, status.Error(codes.Internal, "failed to create shop")
	}
	if exists {
		s.logger.WarnContext(ctx, "slug already exists", slog.String("slug", req.Slug))
		return nil, status.Error(codes.AlreadyExists, "shop slug already exists")
	}

	// Create new shop DTO
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

	// 6️⃣ Save to repository
	id, err := s.repo.CreateShop(ctx, shop)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to create shop in repository",
			slog.String("owner_id", ownerID),
			slog.String("slug", req.Slug),
			slog.String("error", err.Error()),
		)
		return nil, status.Error(codes.Internal, "failed to create shop")
	}

	// 7️⃣ Return response
	return &shoppb.CreateShopResponse{
		ShopId:    id,
		CreatedAt: timestamppb.New(now),
	}, nil
}

func (s *ShopService) ListOwnedShops(ctx context.Context, req *shoppb.ListOwnedShopsRequest) (*shoppb.ListShopsResponse, error) {
	ownerID, userErr := reqCtx.MustGetUserID(ctx)
	if userErr != nil {
		return nil, userErr
	}
	fmt.Printf("ownerID: %v", ownerID)

	shops, err := s.repo.ListByShopOwner(ctx, ownerID)
	fmt.Printf("shops: %v", shops)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to retrieve shops")
	}

	if len(shops) == 0 {
		return nil, status.Error(codes.NotFound, "shop not found")
	}

	return toShopsResponse(shops), nil
}

func (s *ShopService) UpdateShop(ctx context.Context, req *shoppb.UpdateShopRequest) (*shoppb.ShopResponse, error) {
	ownerID, err := reqCtx.MustGetUserID(ctx)

	if err != nil {
		return nil, err
	}

	shop := &dto.ShopDTO{
		OwnerID:     ownerID,
		ShopID:      req.ShopId,
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

func (s *ShopService) DeleteShop(ctx context.Context, req *shoppb.DeleteShopRequest) (*shoppb.DeleteShopResponse, error) {
	ownerID, err := reqCtx.MustGetUserID(ctx)

	if err != nil {
		return nil, err
	}

	affected, err := s.repo.DeleteByOwnerID(ctx, ownerID, req.ShopId)
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
		Name:        s.Name,
		Slug:        s.Slug,
		Description: ptrOrEmpty(s.Description),
		Logo:        ptrOrEmpty(s.Logo),
		IsActive:    s.IsActive,
		CreatedAt:   timestamppb.New(s.CreatedAt),
		UpdatedAt:   timestamppb.New(s.UpdatedAt),
	}
}

func toShopsResponse(shops []*dto.ShopDTO) *shoppb.ListShopsResponse {
	resp := &shoppb.ListShopsResponse{
		Shops: make([]*shoppb.ShopResponse, 0, len(shops)),
	}

	for _, shop := range shops {
		resp.Shops = append(resp.Shops, &shoppb.ShopResponse{
			Id:          shop.ID,
			Name:        shop.Name,
			Slug:        shop.Slug,
			Description: *shop.Description,
			Logo:        *shop.Logo,
			IsActive:    shop.IsActive,
		})
	}

	return resp
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
