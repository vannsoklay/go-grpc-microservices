package proto

import (
	"productservice/internal/domain"
	"productservice/proto/productpb"

	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func MapProductToProto(p *domain.Product) *productpb.Product {

	return &productpb.Product{
		Id:          p.ID,
		ShopId:      p.ShopID,
		Name:        p.Name,
		Category:    p.Category,
		Price:       p.Price,
		Description: nullableString(p.Description),
		Detail:      nullableString(p.Detail),
		CreatedAt:   timestamppb.New(p.CreatedAt),
		UpdatedAt:   timestamppb.New(p.UpdatedAt),
	}
}

func nullableString(s *string) *wrapperspb.StringValue {
	if s == nil || *s == "" {
		return wrapperspb.String("")
	}
	return wrapperspb.String(*s)
}
