package middleware

import (
	"gateway/cache"

	"github.com/gofiber/fiber/v3"
)

func RoleMiddleware(allowedRoles ...string) fiber.Handler {
	return func(c fiber.Ctx) error {
		authCache, ok := c.Locals("authCache").(*cache.AuthCache)
		if !ok {
			return fiber.ErrUnauthorized
		}

		for _, role := range allowedRoles {
			if authCache.Role == role {
				return c.Next()
			}
		}

		return fiber.ErrForbidden
	}
}
