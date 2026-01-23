package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
)

type AuditLog struct {
	Timestamp   time.Time              `json:"timestamp"`
	UserID      string                 `json:"user_id"`
	Username    string                 `json:"username"`
	Method      string                 `json:"method"`
	Path        string                 `json:"path"`
	StatusCode  int                    `json:"status_code"`
	IPAddress   string                 `json:"ip_address"`
	UserAgent   string                 `json:"user_agent"`
	Duration    float64                `json:"duration_ms"`
	QueryParams map[string]interface{} `json:"query_params,omitempty"`
	Body        string                 `json:"body,omitempty"`
	Error       string                 `json:"error,omitempty"`
}

type AuditConfig struct {
	RedisClient    *redis.Client
	LogBody        bool
	LogQueryParams bool
	KeyPrefix      string
	ExpiryDays     int64
	SkipPaths      []string // Paths to skip auditing
}

func AuditMiddleware(config AuditConfig) fiber.Handler {
	if config.KeyPrefix == "" {
		config.KeyPrefix = "audit:"
	}
	if config.ExpiryDays == 0 {
		config.ExpiryDays = 30
	}

	skipMap := make(map[string]bool)
	for _, path := range config.SkipPaths {
		skipMap[path] = true
	}

	return func(c fiber.Ctx) error {
		// Skip audit logging for certain paths
		if skipMap[c.Path()] {
			return c.Next()
		}

		start := time.Now()

		// Get user info from context
		userID := ""
		username := ""

		if uid := c.Locals("userID"); uid != nil {
			userID = fmt.Sprintf("%v", uid)
		}
		if uname := c.Locals("username"); uname != nil {
			username = fmt.Sprintf("%v", uname)
		}

		// Capture request body if needed
		var bodyStr string
		if config.LogBody && c.Method() != fiber.MethodGet {
			bodyStr = string(c.Body())
			// Limit body size to 1KB for storage
			if len(bodyStr) > 1024 {
				bodyStr = bodyStr[:1024] + "..."
			}
		}

		// Capture query params
		var queryParams map[string]interface{}
		if config.LogQueryParams {
			queryParams = make(map[string]interface{})
			c.Request().URI().QueryArgs().VisitAll(func(key, value []byte) {
				queryParams[string(key)] = string(value)
			})
		}

		err := c.Next()

		// Create audit log
		auditLog := AuditLog{
			Timestamp:   start,
			UserID:      userID,
			Username:    username,
			Method:      c.Method(),
			Path:        c.Path(),
			StatusCode:  c.Response().StatusCode(),
			IPAddress:   c.IP(),
			UserAgent:   c.Get("User-Agent"),
			Duration:    float64(time.Since(start).Milliseconds()),
			QueryParams: queryParams,
			Body:        bodyStr,
		}

		if err != nil {
			auditLog.Error = err.Error()
		}

		// Store audit log in Redis
		go storeAuditLog(config.RedisClient, config.KeyPrefix, config.ExpiryDays, &auditLog)

		return err
	}
}

func storeAuditLog(redisClient *redis.Client, keyPrefix string, expiryDays int64, log *AuditLog) {
	ctx := context.Background()

	logJSON, err := json.Marshal(log)
	if err != nil {
		fmt.Printf("Error marshaling audit log: %v\n", err)
		return
	}

	// Store with timestamp for easy querying
	key := fmt.Sprintf("%s%d:%s:%s", keyPrefix, log.Timestamp.Unix(), log.Method, log.Path)

	expiry := time.Duration(expiryDays) * 24 * time.Hour
	if err := redisClient.Set(ctx, key, logJSON, expiry).Err(); err != nil {
		fmt.Printf("Error storing audit log: %v\n", err)
	}
}

// GetAuditLogs retrieves audit logs from Redis (helper function)
func GetAuditLogs(redisClient *redis.Client, pattern string) ([]AuditLog, error) {
	ctx := context.Background()
	var logs []AuditLog

	iter := redisClient.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		val, err := redisClient.Get(ctx, iter.Val()).Result()
		if err != nil {
			continue
		}

		var log AuditLog
		if err := json.Unmarshal([]byte(val), &log); err != nil {
			continue
		}
		logs = append(logs, log)
	}

	return logs, iter.Err()
}
