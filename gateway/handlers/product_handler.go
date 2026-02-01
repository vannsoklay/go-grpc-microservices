package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"gateway/cache"
	"gateway/grpc"
	"hpkg/constants/responses"
	"productservice/proto/productpb"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
)

type RequestHandler struct {
	clients *grpc.GRPCClients
	cache   *cache.RedisCache
}

func NewProductHandler(
	clients *grpc.GRPCClients,
	cache *cache.RedisCache,
) *RequestHandler {
	return &RequestHandler{
		clients: clients,
		cache:   cache,
	}
}

func getAuthContext(c fiber.Ctx) (context.Context, *cache.AuthResp, error) {
	ctx, ok := c.Locals("ctx").(context.Context)
	auth, authOk := c.Locals("auth").(*cache.AuthResp)

	if !ok || ctx == nil || !authOk || auth == nil {
		return nil, nil, fiber.ErrUnauthorized
	}

	return ctx, auth, nil
}

func (h *RequestHandler) ListProductsByShop(c fiber.Ctx) error {
	ctx, _, err := getAuthContext(c)
	if err != nil {
		return responses.Error(c, fiber.StatusUnauthorized, responses.ErrUnauthorizedCode)
	}

	shopID, ok := c.Locals("shop_id").(string)
	if !ok || shopID == "" {
		return responses.Error(c, fiber.StatusForbidden, responses.ErrForbiddenCode)
	}

	limitStr := c.Query("limit", "20")
	limit, _ := strconv.ParseInt(limitStr, 10, 32)

	cursor := c.Query("cursor", "")
	sort := c.Query("az", "za", "old", "new")

	// New unified search & filter
	search := c.Query("search", "") // matches title/name
	filter := c.Query("filter", "") // matches category

	cacheKey := fmt.Sprintf(
		"products:shop:%s:cursor:%s:limit:%d:sort:%s:search:%s:filter:%s",
		shopID, cursor, limit, sort, search, filter,
	)

	// Redis cache
	var cached string
	if h.cache != nil {
		cached, _ = h.cache.Get(ctx, cacheKey)
		if cached != "" {
			var result fiber.Map
			if err := json.Unmarshal([]byte(cached), &result); err == nil {
				return responses.Success(c, fiber.StatusOK, result)
			}
		}
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := h.clients.Product.ListProductsByShop(ctx, &productpb.ListProductsByShopRequest{
		ShopId: shopID,
		Limit:  int32(limit),
		Cursor: cursor,
		Search: search,
		Filter: filter,
		Sort:   sort,
	})
	if err != nil {
		return responses.FromGRPC[any](c, err)
	}

	if h.cache != nil {
		if jsonData, err := json.Marshal(resp); err == nil {
			_ = h.cache.Set(ctx, cacheKey, string(jsonData), 30*time.Second)
		}
	}

	return responses.Success(c, fiber.StatusOK, resp)
}

func (h *RequestHandler) GetProductByID(c fiber.Ctx) error {
	ctx, _, err := getAuthContext(c)
	if err != nil {
		return responses.Error(c, fiber.StatusUnauthorized, responses.ErrUnauthorizedCode)
	}

	productID := c.Params("id")
	if productID == "" {
		return responses.Error(c, fiber.StatusBadRequest, responses.ErrBadRequestCode)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := h.clients.Product.GetProductByID(ctx, &productpb.GetProductRequest{
		ProductId: productID,
	})
	if err != nil {
		return err
	}

	return c.JSON(resp.Product)
}

func (h *RequestHandler) CreateProduct(c fiber.Ctx) error {
	ctx, _, err := getAuthContext(c)
	if err != nil {
		return responses.Error(c, fiber.StatusUnauthorized, responses.ErrUnauthorizedCode)
	}

	var req productpb.CreateProductRequest
	if err := c.Bind().Body(&req); err != nil {
		return responses.Error(c, fiber.StatusBadRequest, responses.ErrBadRequestCode)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := h.clients.Product.CreateProduct(ctx, &req)
	if err != nil {
		return responses.FromGRPC[any](c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "product created",
		"product": resp.Product,
	})
}

func (h *RequestHandler) UpdateProduct(c fiber.Ctx) error {
	ctx, _, err := getAuthContext(c)
	if err != nil {
		return responses.Error(c, fiber.StatusUnauthorized, responses.ErrUnauthorizedCode)
	}

	productID := c.Params("id")
	if productID == "" {
		return responses.Error(c, fiber.StatusBadRequest, responses.ErrBadRequestCode)
	}

	var req productpb.UpdateProductRequest
	if err := c.Bind().Body(&req); err != nil {
		return responses.Error(c, fiber.StatusBadRequest, responses.ErrBadRequestCode)
	}

	req.ProductId = productID

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := h.clients.Product.UpdateProduct(ctx, &req)
	if err != nil {
		return responses.FromGRPC[any](c, err)
	}

	return c.JSON(fiber.Map{
		"message": "product updated",
		"product": resp.Product,
	})
}

func (h *RequestHandler) DeleteProduct(c fiber.Ctx) error {
	ctx, _, err := getAuthContext(c)
	if err != nil {
		return responses.Error(c, fiber.StatusUnauthorized, responses.ErrUnauthorizedCode)
	}

	productID := c.Params("id")
	if productID == "" {
		return responses.Error(c, fiber.StatusBadRequest, responses.ErrBadRequestCode)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := h.clients.Product.DeleteProduct(ctx, &productpb.DeleteProductRequest{
		ProductId: productID,
	})

	if err != nil {
		return responses.FromGRPC[any](c, err)
	}

	return responses.FromGRPC[any](c, err, resp)
}
