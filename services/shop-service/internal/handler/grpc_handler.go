package handler

import (
	"context"
	"shopservice/internal/service"
	"shopservice/proto/shoppb"
)

type ShopHandler struct {
	shoppb.UnimplementedShopServiceServer
	svc *service.ShopService
}

func NewShopHandler(svc *service.ShopService) *ShopHandler {
	return &ShopHandler{svc: svc}
}

func (h *ShopHandler) ValidateShop(ctx context.Context, req *shoppb.ValidateShopRequest) (*shoppb.ValidateShopResponse, error) {
	return h.svc.ValidateShop(ctx, req)
}

func (h *ShopHandler) CreateShop(ctx context.Context, req *shoppb.CreateShopRequest) (*shoppb.CreateShopResponse, error) {
	return h.svc.CreateShop(ctx, req)
}

func (h *ShopHandler) ListOwnedShops(ctx context.Context, req *shoppb.ListOwnedShopsRequest) (*shoppb.ListShopsResponse, error) {
	return h.svc.ListOwnedShops(ctx, req)
}

func (h *ShopHandler) UpdateShop(ctx context.Context, req *shoppb.UpdateShopRequest) (*shoppb.ShopResponse, error) {
	return h.svc.UpdateShop(ctx, req)
}

func (h *ShopHandler) DeleteShop(ctx context.Context, req *shoppb.DeleteShopRequest) (*shoppb.DeleteShopResponse, error) {
	return h.svc.DeleteShop(ctx, req)
}
