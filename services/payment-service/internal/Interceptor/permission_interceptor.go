package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func PermissionInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "no metadata")
		}

		role := md.Get("role")
		if len(role) == 0 {
			return nil, status.Error(codes.PermissionDenied, "missing role")
		}

		// if !Allowed(role[0], info.FullMethod) {
		// 	return nil, status.Error(codes.PermissionDenied, "forbidden")
		// }

		return handler(ctx, req)
	}
}
