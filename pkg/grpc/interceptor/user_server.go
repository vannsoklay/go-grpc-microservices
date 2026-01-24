package interceptor

import (
	"context"
	"time"

	ctxkey "hpkg/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// UserUnaryServerInterceptor validates user ID and logs requests
func UserUnaryServerInterceptor(logger Logger) grpc.UnaryServerInterceptor {
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

		// Extract userID
		userIDs := md.Get("x-user-id")
		if len(userIDs) == 0 || userIDs[0] == "" {
			logger.Warn("missing user ID", "method", method)
			return nil, status.Error(codes.Unauthenticated, "user not authenticated")
		}
		userID := userIDs[0]

		// Attach userID to context
		ctx = context.WithValue(ctx, ctxkey.UserIDKey, userID)
		ctx = metadata.NewOutgoingContext(ctx, md)

		logger.Debug("interceptor: user validated",
			"method", method,
			"user_id", userID,
		)

		// Call handler
		resp, err := handler(ctx, req)

		duration := time.Since(startTime)

		if err != nil {
			st, _ := status.FromError(err)
			logger.Error("interceptor: user request failed",
				"method", method,
				"user_id", userID,
				"error", err.Error(),
				"code", st.Code().String(),
				"duration_ms", duration.Milliseconds(),
			)
			return nil, err
		}

		logger.Info("interceptor: user request completed",
			"method", method,
			"user_id", userID,
			"duration_ms", duration.Milliseconds(),
		)

		return resp, nil
	}
}
