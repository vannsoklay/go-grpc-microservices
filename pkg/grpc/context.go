package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func MustGetUserID(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "shop not authenticated")
	}

	return userID, nil
}

func MustGetShopID(ctx context.Context) (string, error) {
	shopID, ok := ctx.Value(ShopIDKey).(string)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "shop not authenticated")
	}

	return shopID, nil
}
