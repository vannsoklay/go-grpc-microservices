package proto

import (
	"productservice/internal/domain"
	"productservice/proto/v1/productpb"
	"productservice/proto/v1/tagpb"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// MapTagToProto converts a domain.Tag to tagpb.Tag
func MapTagToProto(t *domain.Tag) *productpb.Tag {
	if t == nil {
		return nil
	}
	return &productpb.Tag{
		Id:        t.ID,
		ShopId:    t.ShopID,
		Name:      t.Name,
		Slug:      t.Slug,
		CreatedAt: timestamppb.New(t.CreatedAt),
	}
}

// MapTagToStats converts domain tag and product count to tagpb.TagStats
func MapTagToStats(t *domain.Tag, count int32) *tagpb.TagStats {
	return &tagpb.TagStats{
		Tag:          MapTagToProto(t),
		ProductCount: count,
	}
}

// MapTagToDetail converts tag and product list to productpb.TagDetail
func MapTagToDetail(t *domain.Tag, products []*domain.ProductSummary, total int32) *tagpb.TagDetail {
	Products := make([]*productpb.ProductSummary, 0, len(products))
	for _, p := range products {
		Products = append(Products, MapProductSummaryToProto(p))
	}

	return &tagpb.TagDetail{
		Tag:           MapTagToProto(t),
		Products:      Products,
		TotalProducts: total,
	}
}

// MapProductSummaryToProto converts domain.ProductSummary to productpb.ProductSummary
func MapProductSummaryToProto(p *domain.ProductSummary) *productpb.ProductSummary {
	if p == nil {
		return nil
	}
	return &productpb.ProductSummary{
		Id:       p.ID,
		Name:     p.Name,
		Sku:      &p.SKU,
		Price:    p.Price,
		IsActive: p.IsActive,
	}
}
