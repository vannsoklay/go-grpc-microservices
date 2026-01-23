package server

import (
	"context"
	"fmt"

	sharedErr "hpkg/errors"
	"hpkg/grpc"
	ctxkey "hpkg/grpc"
	"productservice/domain"
	"productservice/proto/productpb"
	"productservice/repository"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type ProductService struct {
	productpb.UnimplementedProductServiceServer
	repo repository.PostgresProductRepository
}

func NewProductService(repo repository.PostgresProductRepository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

func MustGetShopID(ctx context.Context) (string, error) {
	shopID, ok := ctx.Value(ctxkey.ShopIDKey).(string)
	fmt.Printf("shop id %v", shopID)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "shop not authenticated")
	}

	return shopID, nil
}

// ---------------------------
// Helper: Map Domain → Proto
// ---------------------------
func mapDomainToProto(p *domain.Product) *productpb.Product {
	var detail *wrapperspb.StringValue
	if p.Detail != nil {
		detail = wrapperspb.String(*p.Detail)
	} else {
		detail = &wrapperspb.StringValue{Value: ""}
	}

	return &productpb.Product{
		Id:        p.ID,
		ShopId:    p.ShopID,
		Name:      p.Name,
		Category:  p.Category,
		Price:     p.Price,
		Detail:    detail,
		CreatedAt: timestamppb.New(p.CreatedAt),
		UpdatedAt: timestamppb.New(p.UpdatedAt),
	}
}

// ---------------------------
// LIST PRODUCTS BY SHOP
// ---------------------------
func (s *ProductService) ListProductsByShop(
	ctx context.Context,
	req *productpb.ListProductsByShopRequest,
) (*productpb.ListProductsByShopResponse, error) {

	var sortColumn string
	var sortDesc bool
	switch req.Sort {
	case "az":
		sortColumn = "name"
		sortDesc = false
	case "za":
		sortColumn = "name"
		sortDesc = true
	case "old":
		sortColumn = "created_at"
		sortDesc = false
	case "new":
		fallthrough
	default:
		sortColumn = "created_at"
		sortDesc = true
	}

	shopID, err := MustGetShopID(ctx)
	if err != nil {
		return nil, err
	}

	fmt.Printf("shop id %v", shopID)

	products, nextCursor, err := s.repo.ListByShopID(ctx, shopID, req.Search, req.Filter, sortColumn, sortDesc, int(req.Limit), req.Cursor)
	if err != nil {
		return nil, sharedErr.ToGRPC(err)
	}

	resp := &productpb.ListProductsByShopResponse{
		Products: make([]*productpb.Product, 0),
	}

	for _, p := range products {
		resp.Products = append(resp.Products, mapDomainToProto(p))
	}

	if nextCursor != "" {
		resp.NextCursor = wrapperspb.String(nextCursor)
	} else {
		resp.NextCursor = nil // returns null in JSON/gRPC
	}

	return resp, nil
}

// ---------------------------
// GET PRODUCT BY ID
// ---------------------------
func (s *ProductService) GetProductByID(
	ctx context.Context,
	req *productpb.GetProductRequest,
) (*productpb.GetProductResponse, error) {

	// Retrieve shop_id from context (child context works too)
	shopID, err := MustGetShopID(ctx)
	if err != nil {
		return nil, err
	}

	fmt.Printf("shop id %v", shopID)

	product, err := s.repo.GetByID(ctx, req.ProductId)
	if err != nil {
		return nil, sharedErr.ToGRPC(err)
	}

	return &productpb.GetProductResponse{
		Product: mapDomainToProto(product),
	}, nil
}

// ---------------------------
// CREATE PRODUCT
// ---------------------------
func (s *ProductService) CreateProduct(
	ctx context.Context,
	req *productpb.CreateProductRequest,
) (*productpb.CreateProductResponse, error) {

	userID, err := grpc.MustGetUserID(ctx)
	if err != nil {
		return nil, sharedErr.ToGRPC(err)
	}

	fmt.Printf("Creating product for userID: %v\n", userID)

	product, err := s.repo.Create(ctx, domain.CreateProductRequest{
		ShopID:      "6671d2a0-c585-4410-af40-cea0bb305ff8",
		OwnerID:     userID,
		Name:        req.Name,
		Description: req.Description,
		Category:    "category-new1",
		Price:       req.Price,
		Detail:      req.Description,
	})
	if err != nil {
		return nil, sharedErr.ToGRPC(err)
	}

	return &productpb.CreateProductResponse{
		Product: mapDomainToProto(product),
	}, nil
}

// ---------------------------
// UPDATE PRODUCT
// ---------------------------
func (s *ProductService) UpdateProduct(
	ctx context.Context,
	req *productpb.UpdateProductRequest,
) (*productpb.UpdateProductResponse, error) {

	// Map gRPC request to domain update
	updateReq := domain.UpdateProductRequest{
		ID:       req.ProductId,
		Name:     req.Name,
		Category: req.Category,
		Price:    req.Price,
		Detail:   "detail update",
	}

	product, err := s.repo.Update(ctx, updateReq)
	if err != nil {
		return nil, sharedErr.ToGRPC(err)
	}

	return &productpb.UpdateProductResponse{
		Product: mapDomainToProto(product),
	}, nil
}

// ---------------------------
// DELETE PRODUCT
// ---------------------------
func (s *ProductService) DeleteProduct(
	ctx context.Context,
	req *productpb.DeleteProductRequest,
) (*productpb.DeleteProductResponse, error) {

	if err := s.repo.Delete(ctx, req.ProductId); err != nil {
		return nil, sharedErr.ToGRPC(err)
	}

	return &productpb.DeleteProductResponse{
		Message: "product deleted successfully",
	}, nil
}
