package proto

import (
	domain "userservice/internal/domain/entities"
	"userservice/proto/userpb"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func MapUserToProto(u *domain.User) *userpb.UserDetailResponse {
	resp := &userpb.UserDetailResponse{
		Id:         u.ID.String(),
		Name:       u.Name,
		Username:   u.Username,
		Email:      u.Email,
		Bio:        stringPtr(u.Bio),
		IsVerified: u.IsVerified,
		Status:     u.Status,
		CreatedAt:  timestamppb.New(u.CreatedAt),
		UpdatedAt:  timestamppb.New(u.UpdatedAt),
	}

	if u.EmailVerifiedAt != nil {
		resp.EmailVerifiedAt = timestamppb.New(*u.EmailVerifiedAt)
	}
	if u.LastLogin != nil {
		resp.LastLogin = timestamppb.New(*u.LastLogin)
	}
	if u.RoleID != nil {
		resp.RoleId = u.RoleID.String()
	}

	return resp
}

func stringPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
