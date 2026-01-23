package grpc

import (
	"context"
	"time"

	productpb "productservice/proto/productpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProductClient struct {
	client productpb.ProductServiceClient
}

func NewProductClient(addr string) (*ProductClient, error) {
	conn, err := grpc.Dial(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &ProductClient{
		client: productpb.NewProductServiceClient(conn),
	}, nil
}

func (p *ProductClient) GetProduct(ctx context.Context, productID string) (*productpb.GetProductResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	return p.client.GetProductByID(ctx, &productpb.GetProductRequest{
		ProductId: productID,
	})
}
