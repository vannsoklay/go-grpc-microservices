package service

import (
	"context"
	"database/sql"

	"paymentservice/internal/domain"
)

type PaymentService struct {
	db *sql.DB
}

func NewPaymentService(db *sql.DB) *PaymentService {
	return &PaymentService{db: db}
}

func (r *PaymentService) Create(
	ctx context.Context,
	p *domain.Payment,
) (*domain.Payment, error) {

	query := `
		INSERT INTO payments (
			order_id, user_id, amount, currency,
			payment_method, status,
			transaction_id, reference_id,
			processing_fee
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		p.OrderID,
		p.UserID,
		p.Amount,
		p.Currency,
		p.PaymentMethod,
		p.Status,
		p.TransactionID,
		p.ReferenceID,
		p.ProcessingFee,
	).Scan(
		&p.ID,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	return p, err
}

func (r *PaymentService) GetByID(
	ctx context.Context,
	id int32,
) (*domain.Payment, error) {

	query := `
		SELECT
			id, order_id, user_id,
			amount, currency,
			payment_method, status,
			transaction_id, reference_id,
			processing_fee,
			created_at, updated_at
		FROM payments
		WHERE id = $1
	`

	var p domain.Payment
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID,
		&p.OrderID,
		&p.UserID,
		&p.Amount,
		&p.Currency,
		&p.PaymentMethod,
		&p.Status,
		&p.TransactionID,
		&p.ReferenceID,
		&p.ProcessingFee,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	return &p, err
}

func (r *PaymentService) UpdateStatus(
	ctx context.Context,
	id int32,
	status string,
) error {

	_, err := r.db.ExecContext(
		ctx,
		`UPDATE payments SET status=$1, updated_at=now() WHERE id=$2`,
		status,
		id,
	)
	return err
}
