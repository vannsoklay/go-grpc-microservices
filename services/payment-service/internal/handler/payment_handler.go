package handler

import (
	"context"
	"paymentservice/internal/domain"
	"paymentservice/internal/repository"
	"paymentservice/internal/service"
	paymentpb "paymentservice/proto/paymentpb"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PaymentHandler struct {
	paymentpb.UnimplementedPaymentServiceServer
	repo repository.PaymentRepository
	svc  *service.PaymentService
}

func NewPaymentHandler(
	svc *service.PaymentService,
) *PaymentHandler {
	return &PaymentHandler{
		svc: svc,
	}
}

func (h *PaymentHandler) ProcessPayment(
	ctx context.Context,
	req *paymentpb.ProcessPaymentRequest,
) (*paymentpb.ProcessPaymentResponse, error) {

	if h.svc == nil {
		return nil, status.Error(codes.Internal, "payment service not initialized")
	}

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	if req.PaymentMethod == "credit_card" && req.Card == nil {
		return nil, status.Error(
			codes.InvalidArgument,
			"card details are required",
		)
	}
	var p = &domain.Payment{
		OrderID:       req.OrderId,
		UserID:        req.UserId,
		Amount:        req.Amount,
		Currency:      req.Currency,
		PaymentMethod: req.PaymentMethod,
	}

	payment, err := h.svc.Create(ctx, p)
	if err != nil {
		return nil, err
	}

	if payment == nil {
		return nil, status.Error(codes.Internal, "payment result is nil")
	}

	transactionID := "txn_" + time.Now().Format("20160102150405")
	return &paymentpb.ProcessPaymentResponse{
		PaymentId:     payment.ID,
		OrderId:       payment.OrderID,
		TransactionId: transactionID,
		Status:        payment.Status,
		Amount:        payment.Amount,
		ProcessingFee: payment.ProcessingFee,
		Message:       "payment processed successfully",
	}, nil
}

func (s *PaymentHandler) GetPayment(
	ctx context.Context,
	req *paymentpb.GetPaymentRequest,
) (*paymentpb.GetPaymentResponse, error) {

	p, err := s.repo.GetByID(ctx, req.PaymentId)
	if err != nil {
		return nil, err
	}

	return &paymentpb.GetPaymentResponse{
		Id:            p.ID,
		OrderId:       p.OrderID,
		UserId:        p.UserID,
		Amount:        p.Amount,
		Currency:      p.Currency,
		PaymentMethod: p.PaymentMethod,
		Status:        p.Status,
		TransactionId: p.TransactionID,
		ReferenceId:   p.ReferenceID,
		ProcessingFee: p.ProcessingFee,
		CreatedAt:     timestamppb.New(p.CreatedAt),
		UpdatedAt:     timestamppb.New(p.UpdatedAt),
	}, nil
}

func (s *PaymentHandler) VerifyPayment(
	ctx context.Context,
	req *paymentpb.VerifyPaymentRequest,
) (*paymentpb.VerifyPaymentResponse, error) {

	return &paymentpb.VerifyPaymentResponse{
		PaymentId:  req.PaymentId,
		IsVerified: true,
		Status:     "completed",
		Message:    "payment verified",
	}, nil
}

func (s *PaymentHandler) ValidatePayment(
	ctx context.Context,
	req *paymentpb.ValidatePaymentRequest,
) (*paymentpb.ValidatePaymentResponse, error) {

	// example logic
	if req.ExpectedAmount <= 0 {
		return &paymentpb.ValidatePaymentResponse{
			IsValid: false,
			Message: "invalid amount",
		}, nil
	}

	return &paymentpb.ValidatePaymentResponse{
		IsValid: true,
		Message: "payment is valid",
	}, nil
}
