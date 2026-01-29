package grpc

import (
	"context"

	err "hpkg/constants/responses"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func MustGetUserID(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return "", status.Errorf(codes.Unauthenticated, "%s:%s", err.ErrUnauthorizedCode, err.ErrUnauthorizedMsg)
	}

	return userID, nil
}

func MustGetShopID(ctx context.Context) (string, error) {
	shopID, ok := ctx.Value(ShopIDKey).(string)
	if !ok {
		return "", status.Errorf(codes.Unauthenticated, "%s:%s", err.ErrUnauthorizedCode, err.ErrUnauthorizedMsg)
	}

	return shopID, nil
}
