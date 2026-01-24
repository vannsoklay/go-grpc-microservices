package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type AuthCache struct {
	rdb *RedisCache
	ttl time.Duration
}

type AuthResp struct {
	UserID      string   `json:"user_id"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
}

func NewAuthCache(rdb *RedisCache, ttl time.Duration) *AuthCache {
	return &AuthCache{rdb: rdb, ttl: ttl}
}

func (r *AuthCache) GetAuth(ctx context.Context, token string) (*AuthResp, error) {
	val, err := r.rdb.Get(ctx, token)
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var cache AuthResp
	if err := json.Unmarshal([]byte(val), &cache); err != nil {
		return nil, err
	}

	return &cache, nil
}

func (r *AuthCache) SetAuth(
	ctx context.Context,
	token string,
	value *AuthResp,
	ttl time.Duration,
) error {
	data, _ := json.Marshal(value)
	return r.rdb.Set(ctx, token, string(data), ttl)
}
