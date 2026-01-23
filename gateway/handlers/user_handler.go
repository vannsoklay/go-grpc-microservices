package handler

import (
	"context"
	"gateway/cache"
	"gateway/grpc"

	"hpkg/constants"
	"hpkg/constants/response"
	grpcmw "hpkg/grpc"

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
	// --- defensive check ---
	ctx, ok := c.Locals("ctx").(context.Context)
	auth, authOk := c.Locals("auth").(*cache.AuthCache)

	if !ok || ctx == nil || !authOk || auth == nil {
		return response.Error(c, fiber.StatusUnauthorized, constants.ErrUnauthorizedCode)
	}

	// --- gRPC call with safe context ---
	resp, err := h.clients.User.GetUserDetail(ctx, &emptypb.Empty{})
	if err != nil {
		httpErr := grpcmw.ToGRPC(err)
		return response.Error(c, httpErr.Status, httpErr.Code)
	}

	return response.Success(c, fiber.StatusOK, resp)
}
