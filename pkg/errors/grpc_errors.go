package errors

import (
	"context"
	"database/sql"
	stdErrors "errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrNotFound       = stdErrors.New("not found")
	ErrUnauthorized   = stdErrors.New("unauthorized")
	ErrForbidden      = stdErrors.New("forbidden")
	ErrInvalidInput   = stdErrors.New("invalid input")
	ErrInternalServer = stdErrors.New("internal server error")
)

// ToGRPC converts domain errors to gRPC status errors
func ToGRPC(err error) error {
	if err == nil {
		return nil
	}

	// 1️⃣ If already a gRPC status error → return as-is
	if _, ok := status.FromError(err); ok {
		return err
	}

	// 2️⃣ Context errors
	switch {
	case stdErrors.Is(err, context.Canceled):
		return status.Error(codes.Canceled, "request canceled")

	case stdErrors.Is(err, context.DeadlineExceeded):
		return status.Error(codes.DeadlineExceeded, "request timeout")
	}

	// 3️⃣ Infrastructure errors
	if stdErrors.Is(err, sql.ErrNoRows) {
		return status.Error(codes.NotFound, ErrNotFound.Error())
	}

	// 4️⃣ Domain errors
	switch {
	case stdErrors.Is(err, ErrNotFound):
		return status.Error(codes.NotFound, ErrNotFound.Error())

	case stdErrors.Is(err, ErrUnauthorized):
		return status.Error(codes.Unauthenticated, ErrUnauthorized.Error())

	case stdErrors.Is(err, ErrForbidden):
		return status.Error(codes.PermissionDenied, ErrForbidden.Error())

	case stdErrors.Is(err, ErrInvalidInput):
		return status.Error(codes.InvalidArgument, ErrInvalidInput.Error())

	case stdErrors.Is(err, ErrInternalServer):
		return status.Error(codes.Internal, ErrInternalServer.Error())
	}

	// 5️⃣ Fallback (never leak internal error details)
	return status.Error(codes.Internal, "internal server error")
}
