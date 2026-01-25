package grpc

import (
	"authservice/proto/authpb"
	"fmt"
	"paymentservice/proto/paymentpb"
	"productservice/proto/productpb"
	"shopservice/proto/shoppb"
	"userservice/proto/userpb"

	"gateway/grpc/interceptor"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClients struct {
	Auth    authpb.AuthServiceClient
	User    userpb.UserServiceClient
	Product productpb.ProductServiceClient
	Payment paymentpb.PaymentServiceClient
	Shop    shoppb.ShopServiceClient
}

// NewGRPCClients initializes all gRPC clients with connection pooling
func NewGRPCClients() (*GRPCClients, error) {
	clients := &GRPCClients{}

	// Auth Service
	authConn, err := grpc.Dial(":50050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service: %v", err)
	}
	clients.Auth = authpb.NewAuthServiceClient(authConn)

	// User Service
	userConn, err := grpc.Dial(":50055", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithUnaryInterceptor(
		interceptor.UserMetadataUnaryInterceptor(),
	))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service: %v", err)
	}
	clients.User = userpb.NewUserServiceClient(userConn)

	// Shop Service
	shopConn, err := grpc.Dial(":50058", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithChainUnaryInterceptor(
		interceptor.UserMetadataUnaryInterceptor(),
	))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service: %v", err)
	}
	clients.Shop = shoppb.NewShopServiceClient(shopConn)

	// Product Service
	productConn, err := grpc.Dial(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithChainUnaryInterceptor(
		interceptor.UserMetadataUnaryInterceptor(),
		interceptor.ShopMetadataUnaryClientInterceptor(),
	))

	if err != nil {
		return nil, fmt.Errorf("failed to connect to product service: %v", err)
	}
	clients.Product = productpb.NewProductServiceClient(productConn)

	// // Payment Service
	paymentConn, err := grpc.Dial(":50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to payment service: %v", err)
	}
	clients.Payment = paymentpb.NewPaymentServiceClient(paymentConn)

	return clients, nil
}
