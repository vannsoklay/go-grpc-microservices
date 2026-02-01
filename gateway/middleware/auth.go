package middleware

import (
	"authservice/proto/authpb"
	"context"
	"strings"
	"time"

	"gateway/cache"
	"gateway/grpc"

	errors "hpkg/constants/responses"
	ctxkey "hpkg/grpc"

	"github.com/gofiber/fiber/v3"
)

func AuthMiddleware(client *grpc.GRPCClients, authCache *cache.AuthCache) fiber.Handler {
	return func(c fiber.Ctx) error {
		header := c.Get("Authorization")
		if header == "" {
			return errors.Error(c, fiber.StatusUnauthorized, errors.ErrAuthHeaderMissingCode, errors.ErrAuthHeaderMissingMsg)
		}

		token := strings.TrimPrefix(header, "Bearer ")
		if token == "" {
			return errors.Error(c, fiber.StatusUnauthorized, errors.TokenInvalidCode, errors.TokenInvalidMsg)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		var authResp *cache.AuthResp

		// Try Redis
		cached, err := authCache.GetAuth(ctx, token)
		if err == nil && cached != nil {
			authResp = cached
		} else {
			// Validate token via Auth Service
			resp, err := client.Auth.Validate(ctx, &authpb.TokenReq{
				Token: token,
			})

			if err != nil {
				return errors.FromError(c, err)
			}

			authResp = &cache.AuthResp{
				UserID:      resp.UserId,
				Role:        resp.Role,
				Permissions: resp.Permissions,
			}

			_ = authCache.SetAuth(ctx, token, authResp, 10*time.Minute)
		}

		// Attach auth to request context for gRPC
		reqCtx := context.WithValue(context.Background(), ctxkey.UserIDKey, authResp)

		// Save context and auth for handler
		c.Locals("ctx", reqCtx)
		c.Locals("auth", authResp)

		return c.Next()
	}
}
