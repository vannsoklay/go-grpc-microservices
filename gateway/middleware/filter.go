package middleware

import (
	"errors"

	"github.com/gofiber/fiber/v3"

	"hpkg/constants/responses"
)

// ResponseFilter ensures all responses have a consistent structure
func ResponseFilter() fiber.Handler {
	return func(c fiber.Ctx) error {
		// Recover from panics
		defer func() {
			if r := recover(); r != nil {
				_ = responses.Error(c, fiber.StatusInternalServerError, "ERR_INTERNAL")
			}
		}()

		// Execute next middleware/handler
		err := c.Next()
		if err == nil {
			// no error, proceed
			return nil
		}

		// If the error is a grpc.HTTPError (already converted from gRPC)
		var httpErr responses.HTTPError
		if errors.As(err, &httpErr) {
			return responses.Error(c, httpErr.Status, httpErr.Code)
		}

		// Fallback for unknown errors
		return responses.Error(c, fiber.StatusInternalServerError, "ERR_INTERNAL")
	}
}
