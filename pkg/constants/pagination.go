package constants

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ProductCursor struct {
	CreatedAt time.Time
	ID        string
}

func EncodeCursor(c ProductCursor) (string, error) {
	raw := fmt.Sprintf("%s|%s", c.CreatedAt.UTC().Format(time.RFC3339Nano), c.ID)
	return base64.StdEncoding.EncodeToString([]byte(raw)), nil
}

func DecodeCursor(cursor string) (*ProductCursor, error) {
	b, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return nil, fmt.Errorf("invalid cursor encoding")
	}

	parts := strings.Split(string(b), "|")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid cursor format")
	}

	t, err := time.Parse(time.RFC3339Nano, parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid cursor time")
	}

	if _, err := uuid.Parse(parts[1]); err != nil {
		return nil, fmt.Errorf("invalid cursor id")
	}

	return &ProductCursor{
		CreatedAt: t,
		ID:        parts[1],
	}, nil
}
