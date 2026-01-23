package handler

import (
	"context"
	"gateway/grpc"
	paymentpb "paymentservice/proto/paymentpb"
	"time"

	"github.com/gofiber/fiber/v3"
)

type PaymentHandler struct {
	clients *grpc.GRPCClients
}

func NewPaymentHandler(clients *grpc.GRPCClients) *PaymentHandler {
	return &PaymentHandler{clients: clients}
}

// ProcessPayment endpoint
func (h *PaymentHandler) ProcessPayment(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := new(paymentpb.ProcessPaymentRequest)
	if err := c.Bind().Body(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	resp, err := h.clients.Payment.ProcessPayment(ctx, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(resp)
}

// GetPayment endpoint
func (h *PaymentHandler) GetPayment(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := new(paymentpb.GetPaymentRequest)
	if err := c.Bind().Body(req); err != nil { // can also use BodyParser
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	resp, err := h.clients.Payment.GetPayment(ctx, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(resp)
}

// VerifyPayment endpoint
func (h *PaymentHandler) VerifyPayment(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := new(paymentpb.VerifyPaymentRequest)
	if err := c.Bind().Body(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	resp, err := h.clients.Payment.VerifyPayment(ctx, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(resp)
}

// ValidatePayment endpoint
func (h *PaymentHandler) ValidatePayment(c fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := new(paymentpb.ValidatePaymentRequest)
	if err := c.Bind().Body(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	resp, err := h.clients.Payment.ValidatePayment(ctx, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(resp)
}
