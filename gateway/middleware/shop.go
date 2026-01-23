package middleware

import (
	"context"
	"fmt"
	"gateway/cache"
	"shopservice/proto/shoppb"

	ctxkey "hpkg/grpc"

	"hpkg/constants"
	"hpkg/constants/response"

	"github.com/gofiber/fiber/v3"
	"google.golang.org/grpc/metadata"
)

func ShopMiddleware(
	shopClient shoppb.ShopServiceClient,
	shopCache *cache.ShopCache,
) fiber.Handler {
	return func(c fiber.Ctx) error {

		// --- 1. Get context ---
		ctx, ok := c.Locals("ctx").(context.Context)
		if !ok || ctx == nil {
			return response.Error(c, fiber.StatusUnauthorized, constants.ErrUnauthorizedCode, constants.ErrUnauthorizedMsg)
		}

		// --- 2. Get auth info ---
		auth, ok := c.Locals("auth").(*cache.AuthCache)
		if !ok || auth == nil {
			return response.Error(c, fiber.StatusUnauthorized, constants.ErrUnauthorizedCode, constants.ErrUnauthorizedMsg)
		}

		// --- 3. Get shop ID from header ---
		shopID := c.Get("X-Shop-Id")
		if shopID == "" {
			return response.Error(c, fiber.StatusBadRequest, constants.ShopRequiredCode, constants.ErrUnauthorizedMsg)
		}

		// --- 4. Check Redis cache (userID + shopID) ---
		cacheKey := auth.UserID + ":" + shopID
		if _, ok := shopCache.Get(ctx, cacheKey); ok {
			ctx = AttachShopToCtx(ctx, shopID)
			c.Locals("ctx", ctx)
			c.Locals("shop_id", shopID)
			return c.Next()
		}

		// --- 5. Validate shop via shop-service ---
		shopResp, err := shopClient.GetMyShop(ctx, &shoppb.GetMyShopRequest{})
		if err != nil {
			return response.Error(c, fiber.StatusForbidden, constants.ShopAccessDeniedCode, "Cannot access this shop")
		}
		if shopResp.Id != shopID {
			return response.Error(c, fiber.StatusForbidden, constants.ShopAccessDeniedCode, "Shop ID mismatch")
		}

		// --- 6. Cache valid mapping ---
		shopCache.Set(ctx, cacheKey, "1")

		// --- 7. Attach shop ID to context for downstream services ---
		ctx = AttachShopToCtx(ctx, shopID)
		c.Locals("ctx", ctx)
		c.Locals("shop_id", shopID)

		return c.Next()
	}
}

func AttachShopToCtx(ctx context.Context, shopID string) context.Context {
	fmt.Printf("Attaching shopID %v to context\n", shopID)

	// Attach to Go context
	ctx = context.WithValue(ctx, ctxkey.ShopIDKey, shopID)

	// Attach to gRPC outgoing metadata
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}
	md.Set("x-shop-id", shopID)

	return metadata.NewOutgoingContext(ctx, md)
}
