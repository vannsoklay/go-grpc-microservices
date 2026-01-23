package interceptor

import (
	"context"

	ctxkey "hpkg/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func ShopUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		shopIDs := md.Get("x-shop-id")
		if len(shopIDs) == 0 {
			return nil, status.Error(codes.Unauthenticated, "shop not selected")
		}

		ctx = context.WithValue(ctx, ctxkey.ShopIDKey, shopIDs[0])

		return handler(ctx, req)
	}
}
