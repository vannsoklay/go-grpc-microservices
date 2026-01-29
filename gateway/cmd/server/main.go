package main

import (
	"log"

	"github.com/gofiber/fiber/v3"

	"gateway/cache"
	"gateway/grpc"
	"gateway/middleware"
	"gateway/router"
)

func main() {
	app := fiber.New()
	app.Use(middleware.ResponseFilter())

	// Redis
	redisCache := cache.NewRedis(
		"localhost:6379",       // Redis address
		"your_secure_password", // password
		0,                      // DB
	)

	clients, err := grpc.NewGRPCClients()
	if err != nil {
		log.Fatalf("Failed to initialize gRPC clients: %v", err)
	}

	router.Setup(app, clients, redisCache)

	log.Println("API Gateway running on :3000")
	log.Fatal(app.Listen(":3000"))
}
