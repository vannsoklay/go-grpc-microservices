package responses

import (
	"net/http"

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
			Code:   SuccessCode,
		}
	}

	st, ok := status.FromError(err)
	if !ok {
		return HTTPError{
			Status: http.StatusInternalServerError,
			Code:   ErrInternalCode,
		}
	}

	switch st.Code() {
	case codes.InvalidArgument:
		return HTTPError{http.StatusBadRequest, ErrInvalidInputCode}
	case codes.Unauthenticated:
		return HTTPError{http.StatusUnauthorized, ErrUnauthorizedCode}
	case codes.PermissionDenied:
		return HTTPError{http.StatusForbidden, ErrForbiddenCode}
	case codes.NotFound:
		return HTTPError{http.StatusNotFound, ErrNotFoundCode}
	case codes.Unavailable, codes.DeadlineExceeded:
		return HTTPError{http.StatusBadGateway, ErrServiceUnavailableCode}
	default:
		return HTTPError{http.StatusInternalServerError, ErrInternalCode}
	}
}
