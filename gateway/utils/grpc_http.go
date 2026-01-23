package utils

import (
	"github.com/gofiber/fiber/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GRPCToHTTP(c fiber.Ctx, err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"error": "unknown error"})
	}

	switch st.Code() {
	case codes.NotFound:
		return c.Status(404).JSON(fiber.Map{"error": st.Message()})
	case codes.InvalidArgument:
		return c.Status(400).JSON(fiber.Map{"error": st.Message()})
	case codes.Unauthenticated:
		return c.Status(401).JSON(fiber.Map{"error": st.Message()})
	case codes.PermissionDenied:
		return c.Status(403).JSON(fiber.Map{"error": st.Message()})
	default:
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}
}
