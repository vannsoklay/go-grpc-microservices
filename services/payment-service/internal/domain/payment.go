package domain

import (
	"time"
)

type Payment struct {
	ID            string
	OrderID       string
	UserID        string
	Amount        float64
	Currency      string
	PaymentMethod string
	Status        string
	TransactionID string
	ReferenceID   string
	ProcessingFee float64
	ErrorMessage  string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
