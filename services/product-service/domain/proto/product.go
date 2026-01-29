package proto

import (
	"productservice/domain"
	"productservice/proto/productpb"

	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func MapProductToProto(p *domain.Product) *productpb.Product {
	var detail *wrapperspb.StringValue
	if p.Detail != nil {
		detail = wrapperspb.String(*p.Detail)
	} else {
		detail = &wrapperspb.StringValue{Value: ""}
	}

	return &productpb.Product{
		Id:          p.ID,
		ShopId:      p.ShopID,
		Name:        p.Name,
		Category:    p.Category,
		Price:       p.Price,
		Description: p.Description,
		Detail:      detail,
		CreatedAt:   timestamppb.New(p.CreatedAt),
		UpdatedAt:   timestamppb.New(p.UpdatedAt),
	}
}
