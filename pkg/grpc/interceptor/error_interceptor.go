package interceptor

import (
	"context"

	"google.golang.org/grpc"
)

func ErrorUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {

		resp, err := handler(ctx, req)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}
