package cache

import (
	"context"
	"time"
)

type ShopCache struct {
	rdb *RedisCache
	ttl time.Duration
}

func NewShopCache(rdb *RedisCache, ttl time.Duration) *ShopCache {
	return &ShopCache{rdb: rdb, ttl: ttl}
}

func (c *ShopCache) key(userID string) string {
	return "shop:ctx:" + userID
}

func (c *ShopCache) Get(ctx context.Context, userID string) (string, bool) {
	val, err := c.rdb.Get(ctx, c.key(userID))
	if err != nil {
		return "", false
	}
	return val, true
}

func (c *ShopCache) Set(ctx context.Context, userID, shopID string) {
	_ = c.rdb.Set(ctx, c.key(userID), shopID, c.ttl)
}

func (c *ShopCache) Delete(ctx context.Context, userID string) {
	_ = c.rdb.Del(ctx, c.key(userID))
}
