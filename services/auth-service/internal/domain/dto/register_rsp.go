package dto

import "authservice/internal/domain"

type RegisterRsp struct {
	AccessToken  string
	RefreshToken string
	User         domain.UserRsp
}
