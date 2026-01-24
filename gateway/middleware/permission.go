package middleware

import (
	"gateway/cache"
	errors "hpkg/constants/responses"

	"github.com/gofiber/fiber/v3"
)

func PermissionMiddleware(requiredPerms ...string) fiber.Handler {
	return func(c fiber.Ctx) error {
		// Must come after AuthMiddleware
		auth, ok := c.Locals("auth").(*cache.AuthResp)
		if !ok || auth == nil {
			return errors.Error(
				c,
				fiber.StatusUnauthorized,
				errors.ErrUnauthorizedCode,
			)
		}

		if !hasPermissions(auth.Permissions, requiredPerms) {
			return errors.Error(
				c,
				fiber.StatusForbidden,
				errors.ErrForbiddenCode,
			)
		}

		return c.Next()
	}
}
