package server

import (
	"context"
	pkg "hpkg/grpc"
	"productservice/internal/domain"
	"productservice/internal/domain/proto"
	"productservice/internal/repository"
	"productservice/proto/v1/productpb"
	"productservice/proto/v1/tagpb"

	"google.golang.org/protobuf/types/known/emptypb"
)

type TagService struct {
	tagpb.UnimplementedTagServiceServer
	repo repository.PostgresTagRepository
}

func NewTagService(repo repository.PostgresTagRepository) *TagService {
	return &TagService{repo: repo}
}

func (s *TagService) Create(ctx context.Context, req *tagpb.TagCreateRequest) (*tagpb.TagResponse, error) {
	shopID, _ := pkg.MustGetShopID(ctx)
	tag, err := s.repo.Create(ctx, domain.Tag{
		ShopID: shopID,
		Name:   req.Name,
		Slug:   req.GetSlug(),
	})
	if err != nil {
		return nil, err
	}
	return &tagpb.TagResponse{Tag: proto.MapTagToStats(tag, 0)}, nil
}

func (s *TagService) Get(ctx context.Context, req *tagpb.TagGetRequest) (*tagpb.TagResponse, error) {
	shopID, _ := pkg.MustGetShopID(ctx)
	tag, count, err := s.repo.GetByID(ctx, shopID, req.Id, req.IncludeProductCount)
	if err != nil {
		return nil, err
	}
	return &tagpb.TagResponse{Tag: proto.MapTagToStats(tag, count)}, nil
}

func (s *TagService) GetDetail(ctx context.Context, req *tagpb.TagGetDetailRequest) (*tagpb.TagDetailResponse, error) {
	shopID, _ := pkg.MustGetShopID(ctx)
	tag, products, total, err := s.repo.GetDetail(ctx, shopID, req.Id, int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, err
	}
	return &tagpb.TagDetailResponse{
		TagDetail: proto.MapTagToDetail(tag, products, total),
	}, nil
}

func (s *TagService) List(ctx context.Context, req *tagpb.TagListRequest) (*tagpb.TagListResponse, error) {
	shopID, _ := pkg.MustGetShopID(ctx)
	tags, total, err := s.repo.List(ctx, shopID, *req.Search, req.IncludeProductCount, int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, err
	}

	resp := &tagpb.TagListResponse{
		Tags:     make([]*tagpb.TagStats, 0),
		Total:    int32(total),
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	for _, t := range tags {
		resp.Tags = append(resp.Tags, proto.MapTagToStats(&t.Tag, t.ProductCount))
	}
	return resp, nil
}

func (s *TagService) AssignToProduct(ctx context.Context, req *tagpb.TagAssignRequest) (*tagpb.TagAssignResponse, error) {
	shopID, _ := pkg.MustGetShopID(ctx)
	tags, err := s.repo.AssignToProduct(ctx, shopID, req.ProductId, req.TagIds, req.ReplaceExisting)
	if err != nil {
		return nil, err
	}

	resp := &tagpb.TagAssignResponse{AssignedTags: make([]*productpb.Tag, 0)}
	for _, t := range tags {
		resp.AssignedTags = append(resp.AssignedTags, proto.MapTagToProto(t))
	}
	return resp, nil
}

func (s *TagService) RemoveFromProduct(ctx context.Context, req *tagpb.TagRemoveRequest) (*emptypb.Empty, error) {
	shopID, _ := pkg.MustGetShopID(ctx)
	if err := s.repo.RemoveFromProduct(ctx, shopID, req.ProductId, req.TagIds); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *TagService) GetByProduct(ctx context.Context, req *tagpb.TagGetByProductRequest) (*tagpb.TagGetByProductResponse, error) {
	shopID, _ := pkg.MustGetShopID(ctx)
	tags, err := s.repo.GetByProductID(ctx, shopID, req.ProductId)
	if err != nil {
		return nil, err
	}

	resp := &tagpb.TagGetByProductResponse{Tags: make([]*productpb.Tag, 0)}
	for _, t := range tags {
		resp.Tags = append(resp.Tags, proto.MapTagToProto(t))
	}
	return resp, nil
}
