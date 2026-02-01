package main

import (
	"hpkg/db"
	"hpkg/grpc/interceptor"
	"log"
	"log/slog"
	"net"
	"os"

	"productservice/internal/repository"
	service "productservice/internal/service"
	productpb "productservice/proto/productpb"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	db := db.ConnectPostgreSQLDB()
	defer db.Close()

	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	repo := *repository.NewPostgresProductRepository(db, logger)
	productServer := service.NewProductService(repo)

	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.UserUnaryServerInterceptor(logger), interceptor.ShopUnaryServerInterceptor(logger), interceptor.ErrorUnaryInterceptor()))
	productpb.RegisterProductServiceServer(grpcServer, productServer)

	log.Println("Product service listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
