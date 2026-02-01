package server

import (
	"context"
	pkg "hpkg/grpc"
	"productservice/internal/domain"
	"productservice/internal/domain/proto"
	"productservice/internal/repository"
	"productservice/proto/v1/categorypb"
)

type CategoryService struct {
	categorypb.UnimplementedCategoryServiceServer
	repo repository.PostgresCategoryRepository
}

func NewCategoryService(repo repository.PostgresCategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) Create(ctx context.Context, req *categorypb.CreateRequest) (*categorypb.CreateResponse, error) {
	shopID, _ := pkg.MustGetShopID(ctx)
	category, err := s.repo.Create(ctx, domain.Category{
		ShopID:   shopID,
		Name:     req.Name,
		Slug:     req.GetSlug(),
		ParentID: req.GetParentId(),
	})
	if err != nil {
		return nil, err
	}
	return &categorypb.CreateResponse{Category: proto.MapCategoryToProto(category)}, nil
}

func (s *CategoryService) Get(ctx context.Context, req *categorypb.GetRequest) (*categorypb.GetResponse, error) {
	shopID, _ := pkg.MustGetShopID(ctx)
	category, err := s.repo.GetByID(ctx, shopID, req.Id)
	if err != nil {
		return nil, err
	}
	return &categorypb.GetResponse{Category: proto.MapCategoryToProto(category)}, nil
}

func (s *CategoryService) List(ctx context.Context, req *categorypb.ListRequest) (*categorypb.ListResponse, error) {
	shopID, _ := pkg.MustGetShopID(ctx)
	categories, total, err := s.repo.List(ctx, shopID, req.Search, int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, err
	}

	resp := &categorypb.ListResponse{
		Categories: make([]*categorypb.CategoryInfo, 0),
		Total:      int32(total),
	}
	for _, c := range categories {
		resp.Categories = append(resp.Categories, proto.MapCategoryToInfo(c))
	}
	return resp, nil
}

func (s *CategoryService) GetTree(ctx context.Context, req *categorypb.GetTreeRequest) (*categorypb.GetTreeResponse, error) {
	shopID, _ := pkg.MustGetShopID(ctx)
	tree, err := s.repo.GetTree(ctx, shopID)
	if err != nil {
		return nil, err
	}
	return &categorypb.GetTreeResponse{Categories: proto.MapCategoryTreeToProto(tree)}, nil
}
