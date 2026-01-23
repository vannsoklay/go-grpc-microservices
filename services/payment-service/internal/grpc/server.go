package server

import (
	"paymentservice/internal/handler"
	paymentpb "paymentservice/proto/paymentpb"
)

type PaymentServer struct {
	paymentpb.UnimplementedPaymentServiceServer
	handler *handler.PaymentHandler
}

func NewPaymentServer(handler *handler.PaymentHandler) *PaymentServer {
	return &PaymentServer{handler: handler}
}
