package interceptor

import (
	"context"

	ctxkey "hpkg/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func ShopMetadataUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {

		shopID, ok := ctx.Value(ctxkey.ShopIDKey).(string)
		if !ok || shopID == "" {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}

		md.Set("x-shop-id", shopID)

		ctx = metadata.NewOutgoingContext(ctx, md)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
