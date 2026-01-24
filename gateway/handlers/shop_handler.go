package handler

import (
	"context"
	"gateway/cache"
	"gateway/grpc"
	"shopservice/proto/shoppb"

	"hpkg/constants/responses"

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
	auth, authOk := c.Locals("auth").(*cache.AuthResp)

	if !ok || ctx == nil || !authOk || auth == nil {
		return responses.Error(c, fiber.StatusUnauthorized, responses.ErrUnauthorizedCode)
	}

	var req shoppb.CreateShopRequest
	if err := c.Bind().Body(&req); err != nil {
		return responses.Error(c, fiber.StatusBadRequest, responses.ErrInvalidPayloadCode)
	}

	resp, err := h.clients.Shop.CreateShop(ctx, &req)
	if err != nil {
		return responses.FromError(c, err)
	}

	return responses.Success(c, fiber.StatusCreated, resp)
}

func (h *ShopHandler) GetMyShop(c fiber.Ctx) error {
	ctx, ok := c.Locals("ctx").(context.Context)
	auth, authOk := c.Locals("auth").(*cache.AuthResp)

	if !ok || ctx == nil || !authOk || auth == nil {
		return responses.Error(c, fiber.StatusUnauthorized, responses.ErrUnauthorizedCode)
	}

	resp, err := h.clients.Shop.GetMyShop(ctx, &shoppb.GetMyShopRequest{})
	if err != nil {
		return responses.FromError(c, err)
	}

	return responses.Success(c, fiber.StatusOK, resp)
}

func (h *ShopHandler) UpdateShop(c fiber.Ctx) error {
	ctx, ok := c.Locals("ctx").(context.Context)
	auth, authOk := c.Locals("auth").(*cache.AuthResp)

	if !ok || ctx == nil || !authOk || auth == nil {
		return responses.Error(c, fiber.StatusUnauthorized, responses.ErrUnauthorizedCode)
	}

	var req shoppb.UpdateShopRequest
	if err := c.Bind().Body(&req); err != nil {
		return responses.Error(c, fiber.StatusBadRequest, responses.ErrInvalidPayloadCode)
	}

	resp, err := h.clients.Shop.UpdateShop(ctx, &req)
	if err != nil {
		return responses.FromError(c, err)
	}

	return responses.Success(c, fiber.StatusOK, resp)
}

func (h *ShopHandler) DeleteShop(c fiber.Ctx) error {
	ctx, ok := c.Locals("ctx").(context.Context)
	auth, authOk := c.Locals("auth").(*cache.AuthResp)

	if !ok || ctx == nil || !authOk || auth == nil {
		return responses.Error(c, fiber.StatusUnauthorized, responses.ErrUnauthorizedCode)
	}

	resp, err := h.clients.Shop.DeleteShop(ctx, &shoppb.DeleteShopRequest{})
	if err != nil {
		return responses.FromError(c, err)
	}

	return responses.Success(c, fiber.StatusOK, resp)
}
