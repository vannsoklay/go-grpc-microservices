package middleware

import (
	"context"

	ctxkey "hpkg/grpc"

	"google.golang.org/grpc/metadata"
)

// AttachUserMetadata creates a gRPC outgoing context with user ID and roles from Fiber context or cache.
func AttachUserMetadata(ctx context.Context, userID string) context.Context {
	// Attach to Go context
	ctx = context.WithValue(ctx, ctxkey.UserIDKey, userID)

	// Attach to gRPC outgoing metadata
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}
	md.Set("x-user-id", userID)

	return metadata.NewOutgoingContext(ctx, md)
}

func AttachShopMetadata(ctx context.Context, shopID string) context.Context {
	// Attach to Go context
	ctx = context.WithValue(ctx, ctxkey.ShopIDKey, shopID)

	// Attach to gRPC outgoing metadata
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}
	md.Set("x-shop-id", shopID)

	return metadata.NewOutgoingContext(ctx, md)
}
