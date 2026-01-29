package router

import (
	"time"

	"github.com/gofiber/fiber/v3"

	"gateway/cache"
	"gateway/grpc"
	handler "gateway/handlers"
	mdw "gateway/middleware"
	errs "hpkg/constants/responses"
)

func Setup(app *fiber.App, clients *grpc.GRPCClients, redisCache *cache.RedisCache) {

	h := handler.NewAuthHandler(clients)

	app.Post("/register", h.Register)
	app.Post("/api/auth/login", h.Login)

	// api := app.Group("/api")
	// api.Get("/products", mdw.AuthMiddleware(auth, authCache), hp.ListProductsByShop)
	// api.Get("/products/:id", mdw.AuthMiddleware(auth, authCache), hp.GetProductByID)
	// api.Post("/products", mdw.AuthMiddleware(auth, authCache), hp.CreateProduct)
	// api.Put("/products/:id", mdw.AuthMiddleware(auth, authCache), hp.UpdateProduct)
	// api.Delete("/products/:id", mdw.AuthMiddleware(auth, authCache), hp.DeleteProduct)

	RegisterUserRoutes(app, clients, redisCache)
	// register for shop route
	RegisterShopRoutes(app, clients, redisCache)
	// register for product route
	RegisterProductRoutes(app, clients, redisCache)
	// payment route
	// RegisterPaymentRoutes(app, clients, auth, redisCache, redisCache)

	app.Use(func(c fiber.Ctx) error {
		return errs.Error(c, fiber.StatusUnauthorized, errs.ErrNotFoundCode, errs.ErrNotFoundMsg)
	})
}

func RegisterUserRoutes(
	app *fiber.App,
	clients *grpc.GRPCClients,
	redisCache *cache.RedisCache,
) {
	h := handler.NewUserHandler(clients)
	authCache := cache.NewAuthCache(redisCache, 10*time.Minute)

	// Group for /users
	users := app.Group("/api/users")

	// Get user by id (protected)
	users.Get("/me",
		mdw.AuthMiddleware(clients, authCache),
		// mdw.PermissionMiddleware("user:read"),
		h.GetUser,
	)
}

func RegisterShopRoutes(
	app *fiber.App,
	clients *grpc.GRPCClients,
	redisCache *cache.RedisCache,
) {
	h := handler.NewShopHandler(clients)
	authCache := cache.NewAuthCache(redisCache, 10*time.Minute)
	shopCache := cache.NewShopCache(redisCache, 10*time.Minute)

	shops := app.Group("/api/shops")

	// Create shop (rate limited)
	shops.Post(
		"/",
		mdw.AuthMiddleware(clients, authCache),
		// mdw.RateLimitMiddleware(redisCache, mdw.RateLimitConfig{
		// 	MaxRequests: 1,
		// 	WindowSecs:  600, // 10 minutes
		// 	KeyPrefix:   "create_shop",
		// }),
		mdw.PermissionMiddleware("shop:create"),
		h.CreateShop,
	)

	// Get my shop
	shops.Get(
		"/me",
		mdw.AuthMiddleware(clients, authCache),
		mdw.PermissionMiddleware("shop:read"),
		h.ListByShopOwner,
	)

	// Update my shop
	shops.Put(
		"/me",
		mdw.AuthMiddleware(clients, authCache),
		mdw.ShopMiddleware(clients.Shop, shopCache),
		mdw.PermissionMiddleware("shop:update"),
		h.UpdateShop,
	)

	// Delete my shop
	shops.Delete(
		"/me",
		mdw.AuthMiddleware(clients, authCache),
		mdw.ShopMiddleware(clients.Shop, shopCache),
		mdw.PermissionMiddleware("shop:delete"),
		h.DeleteShop,
	)
}

func RegisterProductRoutes(app *fiber.App, clients *grpc.GRPCClients, redisCache *cache.RedisCache) {
	h := handler.NewProductHandler(clients, redisCache)
	authCache := cache.NewAuthCache(redisCache, 10*time.Minute)
	shopCache := cache.NewShopCache(redisCache, 10*time.Minute)

	api := app.Group("/api/products")

	api.Get("", mdw.AuthMiddleware(clients, authCache), mdw.ShopMiddleware(clients.Shop, shopCache), h.ListProductsByShop)
	api.Get("/:id", mdw.AuthMiddleware(clients, authCache), mdw.ShopMiddleware(clients.Shop, shopCache), h.GetProductByID)
	api.Post("", mdw.AuthMiddleware(clients, authCache), mdw.ShopMiddleware(clients.Shop, shopCache), h.CreateProduct)
	api.Put("/:id", mdw.AuthMiddleware(clients, authCache), mdw.ShopMiddleware(clients.Shop, shopCache), h.UpdateProduct)
	api.Delete("/:id", mdw.AuthMiddleware(clients, authCache), mdw.ShopMiddleware(clients.Shop, shopCache), h.DeleteProduct)
}

func RegisterPaymentRoutes(app *fiber.App, clients *grpc.GRPCClients, redisCache *cache.RedisCache) {
	authCache := cache.NewAuthCache(redisCache, 10*time.Minute)
	h := handler.NewPaymentHandler(clients)

	// Group for /payments
	payments := app.Group("/api/payments")

	payments.Post("/", mdw.AuthMiddleware(clients, authCache), mdw.PermissionMiddleware("PermPaymentCreate"), h.ProcessPayment)

	// Get payment info
	payments.Get("/:payment_id", h.GetPayment)

	// Verify payment
	payments.Post("/verify", mdw.AuthMiddleware(clients, authCache), h.VerifyPayment)

	// Validate payment
	payments.Post("/validate", mdw.AuthMiddleware(clients, authCache), h.ValidatePayment)
}
