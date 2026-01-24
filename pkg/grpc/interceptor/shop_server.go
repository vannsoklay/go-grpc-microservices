package interceptor

import (
	"context"
	ctxkey "hpkg/grpc"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// ShopUnaryServerInterceptor validates shop ID and logs requests
func ShopUnaryServerInterceptor(logger Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		startTime := time.Now()
		method := info.FullMethod

		// Extract metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			logger.Warn("missing metadata", "method", method)
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		// Extract shopID
		shopIDs := md.Get("x-shop-id")
		if len(shopIDs) == 0 || shopIDs[0] == "" {
			logger.Warn("missing shop ID", "method", method)
			return nil, status.Error(codes.Unauthenticated, "shop not authenticated")
		}
		shopID := shopIDs[0]

		// Extract userID (optional)
		userIDs := md.Get("x-user-id")
		userID := ""
		if len(userIDs) > 0 {
			userID = userIDs[0]
		}

		// Log incoming request
		logger.Info("interceptor: incoming request",
			"method", method,
			"shop_id", shopID,
			"user_id", userID,
		)

		// Attach shopID and userID to context for child services
		ctx = context.WithValue(ctx, ctxkey.ShopIDKey, shopID)
		if userID != "" {
			ctx = context.WithValue(ctx, ctxkey.UserIDKey, userID)
		}
		ctx = metadata.NewOutgoingContext(ctx, md) // preserve metadata

		// Call the actual handler
		resp, err := handler(ctx, req)

		// Calculate duration
		duration := time.Since(startTime)

		// Log response based on error
		if err != nil {
			st, _ := status.FromError(err)
			logger.Error("interceptor: request failed",
				"method", method,
				"shop_id", shopID,
				"user_id", userID,
				"error", err.Error(),
				"code", st.Code().String(),
				"duration_ms", duration.Milliseconds(),
			)
			return nil, err
		}

		logger.Info("interceptor: request completed",
			"method", method,
			"shop_id", shopID,
			"user_id", userID,
			"duration_ms", duration.Milliseconds(),
		)

		return resp, nil
	}
}
