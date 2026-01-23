package middleware

import (
	"gateway/cache"
	"hpkg/constants"
	"hpkg/constants/response"

	"github.com/gofiber/fiber/v3"
)

func PermissionMiddleware(requiredPerms ...string) fiber.Handler {
	return func(c fiber.Ctx) error {
		// Must come after AuthMiddleware
		auth, ok := c.Locals("auth").(*cache.AuthCache)
		if !ok || auth == nil {
			return response.Error(
				c,
				fiber.StatusUnauthorized,
				constants.ErrUnauthorizedCode,
			)
		}

		if !hasPermissions(auth.Permissions, requiredPerms) {
			return response.Error(
				c,
				fiber.StatusForbidden,
				constants.ErrForbiddenCode,
			)
		}

		return c.Next()
	}
}
