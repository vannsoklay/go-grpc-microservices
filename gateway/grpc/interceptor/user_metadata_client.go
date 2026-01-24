package interceptor

import (
	"context"
	"gateway/cache"
	ctxkey "hpkg/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UserMetadataUnaryInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {

		if ctx == nil {
			ctx = context.Background()
		}

		auth, ok := ctx.Value(ctxkey.UserIDKey).(*cache.AuthResp)
		if !ok || auth == nil {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		md := metadata.New(map[string]string{
			"x-user-id": auth.UserID,
			"x-roles":   auth.Role,
		})

		ctx = metadata.NewOutgoingContext(ctx, md)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
