package main

import (
	"log"
	"net"
	"os"

	"authservice/internal/handler"
	"authservice/internal/service"
	"authservice/proto/authpb"

	"hpkg/db"
	"hpkg/grpc/interceptor"
	auth "hpkg/grpc/middeware"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	listener, err := net.Listen("tcp", ":50050")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	db := db.ConnectPostgreSQLDB()
	defer db.Close()

	// 1.Create ONE grpc server
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(
		interceptor.ErrorUnaryInterceptor(),
	),
	)
	jwtService := &auth.JWTService{
		Secret: os.Getenv("JWT_SECRET"),
	}

	// 2.dependencies
	svc := service.NewAuthService(db, jwtService)
	h := handler.NewAuthHandler(svc)

	// 3.Register service
	authpb.RegisterAuthServiceServer(grpcServer, h)

	//4.Enable reflection (dev only)
	reflection.Register(grpcServer)

	log.Println("Auth Service listening on :50050")

	// 5.Serve the SAME server
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
