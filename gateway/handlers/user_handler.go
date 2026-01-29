package handler

import (
	"context"
	"gateway/grpc"

	"hpkg/constants/responses"

	"github.com/gofiber/fiber/v3"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserHandler struct {
	clients *grpc.GRPCClients
}

func NewUserHandler(clients *grpc.GRPCClients) *UserHandler {
	return &UserHandler{clients: clients}
}

// GetUser endpoint
func (h *UserHandler) GetUser(c fiber.Ctx) error {
	ctx, ok := c.Locals("ctx").(context.Context)

	if !ok || ctx == nil {
		return responses.Error(c, fiber.StatusUnauthorized, responses.ErrUnauthorizedCode)
	}

	resp, err := h.clients.User.GetUserDetail(ctx, &emptypb.Empty{})
	if err != nil {
		return responses.FromError(c, err)
	}

	return responses.Success(c, fiber.StatusOK, resp)
}
