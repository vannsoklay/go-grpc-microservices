package grpc

// contextKey is a private type to avoid collisions with other context keys
type contextKey string

const (
	// ShopIDKey is the key to store/retrieve shop_id from context
	ShopIDKey contextKey = "shop_id"

	// UserIDKey can be used to store user ID if needed
	UserIDKey contextKey = "user_id"
)
