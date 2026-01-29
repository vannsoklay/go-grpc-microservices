package responses

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GRPC(code codes.Code, errCode, errMsg string) error {
	return status.Errorf(code, "%s:%s", errCode, errMsg)
}
