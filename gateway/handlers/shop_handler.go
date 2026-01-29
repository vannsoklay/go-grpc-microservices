package handler

import (
	"context"
	"gateway/grpc"
	"shopservice/proto/shoppb"

	"hpkg/constants/responses"
	pkg "hpkg/grpc"

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

	if !ok || ctx == nil {
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

func (h *ShopHandler) ListByShopOwner(c fiber.Ctx) error {
	ctx, ok := c.Locals("ctx").(context.Context)

	if !ok || ctx == nil {
		return responses.Error(c, fiber.StatusUnauthorized, responses.ErrUnauthorizedCode)
	}

	resp, err := h.clients.Shop.ListOwnedShops(ctx, &shoppb.ListOwnedShopsRequest{})
	if err != nil {
		return responses.FromError(c, err)
	}

	return responses.Success(c, fiber.StatusOK, resp.Shops)
}

func (h *ShopHandler) UpdateShop(c fiber.Ctx) error {
	ctx, ok := c.Locals("ctx").(context.Context)
	shopID, _ := pkg.MustGetShopID(ctx)
	if !ok || ctx == nil {
		return responses.Error(c, fiber.StatusUnauthorized, responses.ErrUnauthorizedCode)
	}

	var req shoppb.UpdateShopRequest
	if err := c.Bind().Body(&req); err != nil {
		return responses.Error(c, fiber.StatusBadRequest, responses.ErrInvalidPayloadCode)
	}
	req.ShopId = shopID

	resp, err := h.clients.Shop.UpdateShop(ctx, &req)
	if err != nil {
		return responses.FromError(c, err)
	}

	return responses.Success(c, fiber.StatusOK, resp)
}

func (h *ShopHandler) DeleteShop(c fiber.Ctx) error {
	ctx, ok := c.Locals("ctx").(context.Context)

	if !ok || ctx == nil {
		return responses.Error(c, fiber.StatusUnauthorized, responses.ErrUnauthorizedCode)
	}
	shopID, _ := pkg.MustGetShopID(ctx)
	resp, err := h.clients.Shop.DeleteShop(ctx, &shoppb.DeleteShopRequest{
		ShopId: shopID,
	})
	if err != nil {
		return responses.FromError(c, err)
	}

	return responses.Success(c, fiber.StatusOK, resp)
}
