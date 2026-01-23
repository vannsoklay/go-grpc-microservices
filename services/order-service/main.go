package main

import (
	"context"
	"fmt"
	"log"

	"orderservice/grpc"
)

func main() {
	client, err := grpc.NewProductClient("localhost:50051")
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	product, err := client.GetProduct(context.Background(), "1969f353-037a-4d60-bc5b-b5c0a51b342c")
	if err != nil {
		log.Fatalf("failed to get product: %v", err)
	}

	fmt.Printf("Product: %+v\n", product)
}
