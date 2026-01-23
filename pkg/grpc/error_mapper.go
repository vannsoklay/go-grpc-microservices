package grpc

import (
	"net/http"

	"hpkg/constants"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type HTTPError struct {
	Status int    `json:"-"`
	Code   string `json:"code"`
}

// Error implements [error].
func (h HTTPError) Error() string {
	return h.Code
}

func ToGRPC(err error) HTTPError {
	if err == nil {
		return HTTPError{
			Status: http.StatusOK,
			Code:   constants.SuccessCode,
		}
	}

	st, ok := status.FromError(err)
	if !ok {
		return HTTPError{
			Status: http.StatusInternalServerError,
			Code:   constants.ErrInternalCode,
		}
	}

	switch st.Code() {
	case codes.InvalidArgument:
		return HTTPError{http.StatusBadRequest, constants.ErrInvalidInputCode}
	case codes.Unauthenticated:
		return HTTPError{http.StatusUnauthorized, constants.ErrUnauthorizedCode}
	case codes.PermissionDenied:
		return HTTPError{http.StatusForbidden, constants.ErrForbiddenCode}
	case codes.NotFound:
		return HTTPError{http.StatusNotFound, constants.ErrNotFoundCode}
	case codes.Unavailable, codes.DeadlineExceeded:
		return HTTPError{http.StatusBadGateway, constants.ErrServiceUnavailableCode}
	default:
		return HTTPError{http.StatusInternalServerError, constants.ErrInternalCode}
	}
}
