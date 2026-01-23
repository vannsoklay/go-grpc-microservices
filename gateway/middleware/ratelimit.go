package middleware

import (
	"fmt"
	"gateway/cache"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
)

type RateLimitConfig struct {
	MaxRequests int
	WindowSecs  int64
	KeyPrefix   string
}

func RateLimitMiddleware(
	redisClient *redis.Client,
	config RateLimitConfig,
) fiber.Handler {

	if config.MaxRequests == 0 {
		config.MaxRequests = 100
	}
	if config.WindowSecs == 0 {
		config.WindowSecs = 60
	}
	if config.KeyPrefix == "" {
		config.KeyPrefix = "rate_limit"
	}

	return func(c fiber.Ctx) error {
		auth, _ := c.Locals("auth").(*cache.AuthCache)

		var key string
		if auth != nil && auth.UserID != "" {
			key = fmt.Sprintf("%s:user:%s", config.KeyPrefix, auth.UserID)
		} else {
			key = fmt.Sprintf("%s:ip:%s", config.KeyPrefix, c.IP())
		}

		ctx := c.Context()

		count, err := redisClient.Incr(ctx, key).Result()
		if err != nil {
			return fiber.NewError(
				fiber.StatusInternalServerError,
				"rate limit check failed",
			)
		}

		if count == 1 {
			redisClient.Expire(
				ctx,
				key,
				time.Duration(config.WindowSecs)*time.Second,
			)
		}

		if count > int64(config.MaxRequests) {
			return fiber.NewError(
				fiber.StatusTooManyRequests,
				"rate limit exceeded",
			)
		}

		c.Set("X-RateLimit-Limit", fmt.Sprintf("%d", config.MaxRequests))
		c.Set(
			"X-RateLimit-Remaining",
			fmt.Sprintf("%d", max(0, config.MaxRequests-int(count))),
		)

		return c.Next()
	}
}
