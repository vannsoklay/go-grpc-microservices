package interceptor

import (
	"context"
	ctxkey "hpkg/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func ShopUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		// Extract shopID
		shopIDs := md.Get("x-shop-id")
		if len(shopIDs) == 0 || shopIDs[0] == "" {
			return nil, status.Error(codes.Unauthenticated, "shop not authenticated")
		}
		shopID := shopIDs[0]

		// Attach shopID to context for child services
		ctx = context.WithValue(ctx, ctxkey.ShopIDKey, shopID)
		ctx = metadata.NewOutgoingContext(ctx, md) // preserve metadata

		return handler(ctx, req)
	}
}
