package interceptor

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
)

// Logger interface for dependency injection
type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

func LoggingUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {

		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)

		if err != nil {
			log.Printf(
				"❌ %s | %v | %v",
				info.FullMethod,
				duration,
				err,
			)
		} else {
			log.Printf(
				"✅ %s | %v",
				info.FullMethod,
				duration,
			)
		}

		return resp, err
	}
}
