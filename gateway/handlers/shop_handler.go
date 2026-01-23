package handler

import (
	"context"
	"gateway/cache"
	"gateway/grpc"
	"shopservice/proto/shoppb"

	"hpkg/constants"
	"hpkg/constants/response"
	grpcmw "hpkg/grpc"

	"github.com/gofiber/fiber/v3"
)

type ShopHandler struct {
	clients *grpc.GRPCClients
}

func NewShopHandler(clients *grpc.GRPCClients) *ShopHandler {
	return &ShopHandler{clients: clients}
}

func (h *ShopHandler) CreateShop(c fiber.Ctx) error {
	ctx, ok := c.Locals("ctx").(context.Context)
	auth, authOk := c.Locals("auth").(*cache.AuthCache)

	if !ok || ctx == nil || !authOk || auth == nil {
		return response.Error(c, fiber.StatusUnauthorized, constants.ErrUnauthorizedCode)
	}

	var req shoppb.CreateShopRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, constants.ErrInvalidPayloadCode)
	}

	resp, err := h.clients.Shop.CreateShop(ctx, &req)
	if err != nil {
		httpErr := grpcmw.ToGRPC(err)
		return response.Error(c, httpErr.Status, httpErr.Code)
	}

	return response.Success(c, fiber.StatusCreated, resp)
}

func (h *ShopHandler) GetMyShop(c fiber.Ctx) error {
	ctx, ok := c.Locals("ctx").(context.Context)
	auth, authOk := c.Locals("auth").(*cache.AuthCache)

	if !ok || ctx == nil || !authOk || auth == nil {
		return response.Error(c, fiber.StatusUnauthorized, constants.ErrUnauthorizedCode)
	}

	resp, err := h.clients.Shop.GetMyShop(ctx, &shoppb.GetMyShopRequest{})
	if err != nil {
		httpErr := grpcmw.ToGRPC(err)
		return response.Error(c, httpErr.Status, httpErr.Code)
	}

	return response.Success(c, fiber.StatusOK, resp)
}

func (h *ShopHandler) UpdateShop(c fiber.Ctx) error {
	ctx, ok := c.Locals("ctx").(context.Context)
	auth, authOk := c.Locals("auth").(*cache.AuthCache)

	if !ok || ctx == nil || !authOk || auth == nil {
		return response.Error(c, fiber.StatusUnauthorized, constants.ErrUnauthorizedCode)
	}

	var req shoppb.UpdateShopRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, constants.ErrInvalidPayloadCode)
	}

	resp, err := h.clients.Shop.UpdateShop(ctx, &req)
	if err != nil {
		httpErr := grpcmw.ToGRPC(err)
		return response.Error(c, httpErr.Status, httpErr.Code)
	}

	return response.Success(c, fiber.StatusOK, resp)
}

func (h *ShopHandler) DeleteShop(c fiber.Ctx) error {
	ctx, ok := c.Locals("ctx").(context.Context)
	auth, authOk := c.Locals("auth").(*cache.AuthCache)

	if !ok || ctx == nil || !authOk || auth == nil {
		return response.Error(c, fiber.StatusUnauthorized, constants.ErrUnauthorizedCode)
	}

	resp, err := h.clients.Shop.DeleteShop(ctx, &shoppb.DeleteShopRequest{})
	if err != nil {
		httpErr := grpcmw.ToGRPC(err)
		return response.Error(c, httpErr.Status, httpErr.Code)
	}

	return response.Success(c, fiber.StatusOK, resp)
}
