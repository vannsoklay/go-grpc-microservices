package handler

import (
	server "gateway/grpc/server"

	"github.com/gofiber/fiber/v3"
)

type AuthHandler struct {
	auth *server.AuthClient
}

func NewAuthHandler(a *server.AuthClient) *AuthHandler {
	return &AuthHandler{auth: a}
}

func (h *AuthHandler) Register(c fiber.Ctx) error {
	var body struct {
		Name     string `json:"name"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.Bind().Body(&body); err != nil {
		return fiber.ErrBadRequest
	}

	resp, err := h.auth.Register(
		body.Name,
		body.Username,
		body.Email,
		body.Password,
	)

	if err != nil {
		return fiber.ErrUnauthorized
	}

	return c.JSON(resp)
}

func (h *AuthHandler) Login(c fiber.Ctx) error {
	var body struct {
		Email    *string `json:"email"`
		Password string  `json:"password"`
	}

	if err := c.Bind().Body(&body); err != nil {
		return fiber.ErrBadRequest
	}

	resp, err := h.auth.Login(
		body.Email,
		body.Password,
	)

	if err != nil {
		return fiber.ErrUnauthorized
	}

	return c.JSON(resp)
}
