package main

import (
	"log"
	"net"

	"userservice/internal/handler"
	"userservice/internal/service"
	"userservice/proto/userpb"

	"hpkg/db"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"hpkg/grpc/interceptor"
)

func main() {
	listener, err := net.Listen("tcp", ":50055")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	db := db.ConnectPostgreSQLDB()
	defer db.Close()

	// reate ONE grpc server
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(interceptor.AuthUnaryServerInterceptor()))

	// dependencies
	svc := service.NewUserService(db)
	h := handler.NewUserHandler(svc)

	// Register service
	userpb.RegisterUserServiceServer(grpcServer, h)

	//Enable reflection (dev only)
	reflection.Register(grpcServer)

	log.Println("User Service listening on :50055")

	// Serve the SAME server
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
