package main

import (
	"hpkg/grpc/interceptor"
	"log"
	"net"
	"paymentservice/internal/db"
	"paymentservice/internal/handler"
	"paymentservice/internal/service"
	"paymentservice/proto/paymentpb"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatal(err)
	}

	dsn := "postgres://admin:admin123@localhost:5432/mydb?sslmode=disable"

	dbConn, err := db.NewPostgres(dsn)
	if err != nil {
		log.Fatal(err)
	}

	svc := service.NewPaymentService(dbConn)
	h := handler.NewPaymentHandler(svc)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.RecoveryUnaryInterceptor(),
			interceptor.LoggingUnaryInterceptor(),
		),
	)
	paymentpb.RegisterPaymentServiceServer(grpcServer, h)

	log.Println("Payment gRPC running on :50053")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
