package server

import (
	"context"
	"fmt"

	pkg "hpkg/grpc"

	"productservice/internal/domain"
	"productservice/internal/domain/proto"
	"productservice/internal/repository"
	"productservice/proto/productpb"

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

	shopID, err := pkg.MustGetShopID(ctx)
	if err != nil {
		return nil, err
	}

	products, nextCursor, totalCount, totalAllCount, err := s.repo.ListByShopID(ctx, shopID, req.Search, req.Filter, sortColumn, sortDesc, int(req.Limit), req.Cursor)
	if err != nil {
		return nil, err
	}

	resp := &productpb.ListProductsByShopResponse{
		Products:      make([]*productpb.Product, 0),
		PageSize:      int32(len(products)),
		TotalCount:    int32(totalCount),
		TotalAllCount: int32(totalAllCount),
	}

	for _, p := range products {
		resp.Products = append(resp.Products, proto.MapProductToProto(p))
	}

	if nextCursor != "" {
		resp.NextCursor = wrapperspb.String(nextCursor)
	} else {
		resp.NextCursor = nil
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
	shopID, err := pkg.MustGetShopID(ctx)
	if err != nil {
		return nil, err
	}

	fmt.Printf("shop id %v", shopID)

	product, err := s.repo.GetByID(ctx, req.ProductId)
	if err != nil {
		return nil, err
	}

	return &productpb.GetProductResponse{
		Product: proto.MapProductToProto(product),
	}, nil
}

// ---------------------------
// CREATE PRODUCT
// ---------------------------
func (s *ProductService) CreateProduct(
	ctx context.Context,
	req *productpb.CreateProductRequest,
) (*productpb.CreateProductResponse, error) {

	userID, _ := pkg.MustGetUserID(ctx)
	shopID, _ := pkg.MustGetShopID(ctx)

	product, err := s.repo.Create(ctx, domain.CreateProductRequest{
		ShopID:      shopID,
		OwnerID:     userID,
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Price:       req.Price,
		Detail:      req.Detail,
	})

	if err != nil {
		return nil, err
	}

	return &productpb.CreateProductResponse{
		Product: proto.MapProductToProto(product),
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
		ID:          req.ProductId,
		Name:        req.Name,
		Category:    req.Category,
		Price:       req.Price,
		Description: req.Description,
		Detail:      req.Detail,
	}

	product, err := s.repo.Update(ctx, updateReq)
	if err != nil {
		return nil, err
	}

	return &productpb.UpdateProductResponse{
		Product: proto.MapProductToProto(product),
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
		return nil, err
	}

	return &productpb.DeleteProductResponse{
		Message: "product deleted successfully",
	}, nil
}
