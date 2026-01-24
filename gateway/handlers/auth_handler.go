package handler

import (
	"authservice/proto/authpb"
	"gateway/grpc"

	"github.com/gofiber/fiber/v3"
)

type AuthHandler struct {
	clients *grpc.GRPCClients
}

func NewAuthHandler(a *grpc.GRPCClients) *AuthHandler {
	return &AuthHandler{clients: a}
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

	resp, err := h.clients.Auth.Register(
		c.Context(),
		&authpb.RegisterReq{
			Name:     body.Name,
			Username: body.Username,
			Email:    body.Email,
			Password: body.Password,
		},
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

	resp, err := h.clients.Auth.Login(
		c.Context(),
		&authpb.LoginReq{
			Email:    body.Email,
			Password: body.Password,
		},
	)

	if err != nil {
		return fiber.ErrUnauthorized
	}

	return c.JSON(resp)
}
