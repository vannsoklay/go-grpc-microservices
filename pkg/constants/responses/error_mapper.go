package responses

import (
	"errors"
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type HTTPError struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Error implements [error].
func (h HTTPError) Error() string {
	return h.Code
}

type ServiceError struct {
	GRPCCode codes.Code
	Code     string // custom code for client
	Message  string
}

func (e *ServiceError) Error() string {
	return e.Message
}

// helper constructors
func NewShopLimitError(max int) error {
	return &ServiceError{
		GRPCCode: codes.FailedPrecondition,
		Code:     "SHOP_LIMIT_EXCEEDED",
		Message:  fmt.Sprintf("You can only create up to %d shops", max),
	}
}

func NewOrderLimitError(max int) error {
	return &ServiceError{
		GRPCCode: codes.FailedPrecondition,
		Code:     "ORDER_LIMIT_EXCEEDED",
		Message:  fmt.Sprintf("You can only create up to %d orders", max),
	}
}

// ToGRPC converts a gRPC error into an HTTPError with custom message
func ToGRPC(err error) HTTPError {
	if err == nil {
		return HTTPError{Status: http.StatusOK, Code: SuccessCode, Message: "success"}
	}

	// 1️⃣ If it's a ServiceError, use custom code/message
	var svcErr *ServiceError
	if errors.As(err, &svcErr) {
		var httpStatus int
		switch svcErr.GRPCCode {
		case codes.InvalidArgument:
			httpStatus = http.StatusBadRequest
		case codes.FailedPrecondition:
			httpStatus = http.StatusBadRequest
		case codes.AlreadyExists:
			httpStatus = http.StatusConflict
		case codes.NotFound:
			httpStatus = http.StatusNotFound
		default:
			httpStatus = http.StatusInternalServerError
		}

		return HTTPError{
			Status:  httpStatus,
			Code:    svcErr.Code,    // dynamic custom code
			Message: svcErr.Message, // dynamic custom message
		}
	}

	// 2️⃣ Otherwise fallback to normal gRPC status
	st, ok := status.FromError(err)
	if !ok {
		return HTTPError{Status: http.StatusInternalServerError, Code: ErrInternalCode, Message: err.Error()}
	}

	return HTTPError{
		Status:  http.StatusInternalServerError,
		Code:    st.Code().String(),
		Message: st.Message(),
	}
}
