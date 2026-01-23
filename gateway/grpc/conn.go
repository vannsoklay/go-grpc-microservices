package grpc

import "google.golang.org/grpc"

func NewConn(addr string) (*grpc.ClientConn, error) {
	return grpc.Dial(addr, grpc.WithInsecure())
}
