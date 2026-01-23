package middleware

import (
	"context"
	"gateway/cache"
	"strings"

	"github.com/gofiber/fiber/v3"
	"google.golang.org/grpc/metadata"
)

// AttachUserMetadata creates a gRPC outgoing context with user ID and roles from Fiber context or cache.
func AttachUserMetadata(c fiber.Ctx, auth *cache.AuthCache, ctx context.Context) context.Context {
	// Get userID from Fiber locals or fallback to cache
	userIDRaw := c.Locals("userID", auth.UserID)
	userID, _ := userIDRaw.(string) // assert to string

	// Get roles from Fiber locals or fallback to cache
	rolesRaw := c.Locals("role", auth.Role)
	var roles []string

	switch v := rolesRaw.(type) {
	case string:
		roles = []string{v}
	case []string:
		roles = v
	default:
		roles = []string{}
	}

	md := metadata.New(map[string]string{
		"x-user-id": userID,
		"x-roles":   strings.Join(roles, ","),
	})

	return metadata.NewOutgoingContext(ctx, md)
}
