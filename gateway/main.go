package main

import (
	"log"

	"github.com/gofiber/fiber/v3"

	"authservice/proto/authpb"
	"gateway/cache"
	"gateway/grpc"
	server "gateway/grpc/server"
	"gateway/middleware"
	"gateway/router"
)

func main() {
	app := fiber.New()
	app.Use(middleware.ResponseFilter())

	// Redis
	authCache := cache.NewAuthRedisCache(
		"localhost:6379",       // Redis address
		"your_secure_password", // password
		0,                      // DB
	)

	redisCache := cache.NewRedis(
		"localhost:6379",       // Redis address
		"your_secure_password", // password
		0,                      // DB
	)

	clients, err := grpc.NewGRPCClients()
	if err != nil {
		log.Fatalf("Failed to initialize gRPC clients: %v", err)
	}

	conn, err := grpc.NewConn("localhost:50050")
	if err != nil {
		log.Fatal(err)
	}
	authClient := server.NewAuthClient(
		authpb.NewAuthServiceClient(conn),
	)

	router.Setup(app, authClient, clients, authCache, redisCache)

	log.Println("API Gateway running on :3000")
	log.Fatal(app.Listen(":3000"))
}
