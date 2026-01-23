package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type AuthRedisCache struct {
	client *redis.Client
}

type AuthCache struct {
	UserID      string   `json:"user_id"`
	Username    string   `json:"username"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
}

func NewAuthRedisCache(addr, password string, db int) *AuthRedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,     // "localhost:6379"
		Password: password, // "" if none
		DB:       db,
	})

	return &AuthRedisCache{client: rdb}
}

func (r *AuthRedisCache) GetAuth(ctx context.Context, token string) (*AuthCache, error) {
	val, err := r.client.Get(ctx, token).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var cache AuthCache
	if err := json.Unmarshal([]byte(val), &cache); err != nil {
		return nil, err
	}

	return &cache, nil
}

func (r *AuthRedisCache) SetAuth(
	ctx context.Context,
	token string,
	value *AuthCache,
	ttl time.Duration,
) error {
	data, _ := json.Marshal(value)
	return r.client.Set(ctx, token, data, ttl).Err()
}
