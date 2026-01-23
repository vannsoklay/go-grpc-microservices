package grpc

import (
	"context"
	"time"

	"authservice/proto/authpb"
)

type AuthClient struct {
	client authpb.AuthServiceClient
}

func NewAuthClient(c authpb.AuthServiceClient) *AuthClient {
	return &AuthClient{client: c}
}

func (a *AuthClient) Register(name string, username string, email string, password string) (*authpb.RegsiterResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return a.client.Register(ctx, &authpb.RegisterReq{
		Name:     name,
		Email:    email,
		Username: username,
		Password: password,
	})
}

func (a *AuthClient) Login(email *string, password string) (*authpb.LoginResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return a.client.Login(ctx, &authpb.LoginReq{
		Email:    email,
		Password: password,
	})
}

func (a *AuthClient) Validate(token string) (*authpb.ValidateTokenResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return a.client.Validate(ctx, &authpb.TokenReq{
		Token: token,
	})
}
