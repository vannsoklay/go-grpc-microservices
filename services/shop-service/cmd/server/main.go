package main

import (
	"log"
	"log/slog"
	"net"
	"os"

	"shopservice/internal/handler"
	"shopservice/internal/infrastructure/persistence"
	"shopservice/internal/service"
	"shopservice/proto/shoppb"

	"hpkg/db"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"hpkg/grpc/interceptor"
)

func main() {
	listener, err := net.Listen("tcp", ":50058")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	db := db.ConnectPostgreSQLDB()
	defer db.Close()

	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	// reate ONE grpc server
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.UserUnaryServerInterceptor(logger), interceptor.ErrorUnaryInterceptor()))

	// dependencies
	repo := persistence.NewPostgresShopRepository(db, logger)
	svc := service.NewShopService(repo, logger)
	h := handler.NewShopHandler(svc)

	// Register service
	shoppb.RegisterShopServiceServer(grpcServer, h)

	//Enable reflection (dev only)
	reflection.Register(grpcServer)

	log.Println("Shop Service listening on :50058")

	// Serve the SAME server
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
