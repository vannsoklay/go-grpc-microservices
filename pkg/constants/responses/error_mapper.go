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

// ToGRPC converts a gRPC error into an HTTPError with custom message
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

	var httpStatus int
	var code string

	switch st.Code() {
	case codes.InvalidArgument:
		httpStatus = http.StatusBadRequest
		code = ErrInvalidInputCode
	case codes.FailedPrecondition:
		httpStatus = http.StatusBadRequest
		code = ShopLimitExceededCode
	case codes.AlreadyExists:
		httpStatus = http.StatusConflict
		code = "SHOP_SLUG_EXISTS"
	case codes.Unauthenticated:
		httpStatus = http.StatusUnauthorized
		code = ErrUnauthorizedCode
	case codes.PermissionDenied:
		httpStatus = http.StatusForbidden
		code = ErrForbiddenCode
	case codes.NotFound:
		httpStatus = http.StatusNotFound
		code = ErrNotFoundCode
	case codes.Unavailable, codes.DeadlineExceeded:
		httpStatus = http.StatusBadGateway
		code = ErrServiceUnavailableCode
	default:
		httpStatus = http.StatusInternalServerError
		code = ErrInternalCode
	}

	return HTTPError{
		Status: httpStatus,
		Code:   code,
	}
}
