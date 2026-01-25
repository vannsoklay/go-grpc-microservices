package middleware

import (
	"context"
	"gateway/cache"
	"shopservice/proto/shoppb"

	"hpkg/constants/responses"

	"github.com/gofiber/fiber/v3"
)

func ShopMiddleware(
	shopClient shoppb.ShopServiceClient,
	shopCache *cache.ShopCache,
) fiber.Handler {
	return func(c fiber.Ctx) error {

		// --- 1. Get context ---
		ctx, ok := c.Locals("ctx").(context.Context)
		if !ok || ctx == nil {
			return responses.Error(c, fiber.StatusUnauthorized, responses.ErrUnauthorizedCode, responses.ErrUnauthorizedMsg)
		}

		// --- 2. Get auth info ---
		auth, ok := c.Locals("auth").(*cache.AuthResp)
		if !ok || auth == nil {
			return responses.Error(c, fiber.StatusUnauthorized, responses.ErrUnauthorizedCode, responses.ErrUnauthorizedMsg)
		}

		// --- 3. Get shop ID from header ---
		shopID := c.Get("X-Shop-Id")
		if shopID == "" {
			return responses.Error(c, fiber.StatusBadRequest, responses.ShopRequiredCode, responses.ErrUnauthorizedMsg)
		}

		// --- 4. Check Redis cache (userID + shopID) ---
		cacheKey := auth.UserID + ":" + shopID
		if _, ok := shopCache.Get(ctx, cacheKey); ok {
			ctx = AttachShopMetadata(ctx, shopID)
			c.Locals("ctx", ctx)
			c.Locals("shop_id", shopID)
			return c.Next()
		}

		// --- 5. Validate shop via shop-service ---
		shopResp, err := shopClient.ValidateShop(ctx, &shoppb.ValidateShopRequest{
			ShopId: shopID,
		})

		if err != nil {
			return responses.Error(c, fiber.StatusForbidden, responses.ShopAccessDeniedCode, "Cannot access this shop")
		}
		if shopResp.Id != shopID {
			return responses.Error(c, fiber.StatusForbidden, responses.ShopAccessDeniedCode, "Shop ID mismatch")
		}

		// --- 6. Cache valid mapping ---
		shopCache.Set(ctx, cacheKey, "1")

		// --- 7. Attach shop ID to context for downstream services ---
		ctx = AttachShopMetadata(ctx, shopID)
		c.Locals("ctx", ctx)
		c.Locals("shop_id", shopID)

		return c.Next()
	}
}
