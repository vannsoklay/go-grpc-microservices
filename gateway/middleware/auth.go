package middleware

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gateway/cache"
	server "gateway/grpc/server"

	"hpkg/constants"
	"hpkg/constants/response"

	"github.com/gofiber/fiber/v3"

	grpcintcp "hpkg/grpc/interceptor"
)

func AuthMiddleware(authClient *server.AuthClient, redis *cache.AuthRedisCache) fiber.Handler {
	return func(c fiber.Ctx) error {
		header := c.Get("Authorization")
		if header == "" {
			return response.Error(c, fiber.StatusUnauthorized, constants.ErrAuthHeaderMissingCode)
		}

		token := strings.TrimPrefix(header, "Bearer ")
		if token == "" {
			return response.Error(c, fiber.StatusUnauthorized, constants.ErrTokenInvalidCode)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		var authCache *cache.AuthCache

		// Try Redis
		cached, err := redis.GetAuth(ctx, token)
		if err == nil && cached != nil {
			authCache = cached
		} else {
			// Validate token via Auth Service
			resp, err := authClient.Validate(token)
			fmt.Printf("err %v", err)
			if err != nil {
				return response.Error(c, fiber.StatusUnauthorized, constants.ErrTokenExpiredCode)
			}

			authCache = &cache.AuthCache{
				UserID:      resp.UserId,
				Role:        resp.Role,
				Permissions: resp.Permissions,
			}

			_ = redis.SetAuth(ctx, token, authCache, 10*time.Minute)
		}

		// Attach auth to request context for gRPC
		reqCtx := context.WithValue(context.Background(), grpcintcp.AuthContextKey, authCache)

		// Save context and auth for handler
		c.Locals("ctx", reqCtx)
		c.Locals("auth", authCache)

		return c.Next()
	}
}
