package responses

import (
	"net/http"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type HTTPError struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
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
			Status:  http.StatusInternalServerError,
			Code:    ErrInternalCode,
			Message: ErrInternalMsg,
		}
	}

	var httpStatus int
	var code string
	var message string

	switch st.Code() {
	case codes.InvalidArgument:
		httpStatus = http.StatusBadRequest
		code = ErrInvalidInputCode
	case codes.FailedPrecondition:
		httpStatus = http.StatusBadRequest
		code = extractMessage(st.Message()).Code
		message = extractMessage(st.Message()).Message
	case codes.AlreadyExists:
		httpStatus = http.StatusConflict
		code = extractMessage(st.Message()).Code
		message = extractMessage(st.Message()).Message
	case codes.Unauthenticated:
		httpStatus = http.StatusUnauthorized
		code = extractMessage(st.Message()).Code
		message = extractMessage(st.Message()).Message
	case codes.PermissionDenied:
		httpStatus = http.StatusForbidden
		code = extractMessage(st.Message()).Code
		message = extractMessage(st.Message()).Message
	case codes.NotFound:
		httpStatus = http.StatusNotFound
		code = extractMessage(st.Message()).Code
		message = extractMessage(st.Message()).Message
	case codes.Unavailable, codes.DeadlineExceeded:
		httpStatus = http.StatusBadGateway
		code = extractMessage(st.Message()).Code
		message = extractMessage(st.Message()).Message
	default:
		httpStatus = http.StatusInternalServerError
		code = extractMessage(st.Message()).Code
		message = extractMessage(st.Message()).Message
	}

	return HTTPError{
		Status:  httpStatus,
		Code:    code,
		Message: message,
	}
}

func extractMessage(errMsg string) ErrorPayload {
	parts := strings.SplitN(errMsg, ":", 2)

	if len(parts) == 2 {
		return ErrorPayload{
			Code:    strings.TrimSpace(parts[0]),
			Message: strings.TrimSpace(parts[1]),
		}
	}

	// fallback if no code prefix
	return ErrorPayload{
		Code:    "",
		Message: errMsg,
	}
}
