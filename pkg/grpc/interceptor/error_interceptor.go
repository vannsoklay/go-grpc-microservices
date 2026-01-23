package interceptor

import (
	"context"
	appErr "hpkg/grpc"

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
			return nil, appErr.ToGRPC(err)
		}

		return resp, nil
	}
}
