package responses

import (
	"fmt"
)

// CustomError implements the error interface
type CustomError struct {
	Code    string
	Message string
	Err     error
	Status  int
}

// Error implements the error interface (required method)
func (e *CustomError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap allows error chain inspection
func (e *CustomError) Unwrap() error {
	return e.Err
}

// New creates a new CustomError
func New(code, message string, status int) error {
	return &CustomError{
		Code:    code,
		Message: message,
		Status:  status,
	}
}

// NewWithError creates a CustomError with wrapped error
func NewWithError(code, message string, status int, err error) error {
	return &CustomError{
		Code:    code,
		Message: message,
		Status:  status,
		Err:     err,
	}
}

// ValidationError creates a validation error
func ValidationServiceError(message string) error {
	return New("ERR_VALIDATION", message, 400)
}

// NotFoundError creates a not found error
func NotFoundServiceError(message string) error {
	return New("ERR_NOT_FOUND", message, 404)
}

// UnauthorizedError creates an unauthorized error
func UnauthorizedServiceError(message string) error {
	return New("ERR_UNAUTHORIZED", message, 401)
}

// ForbiddenError creates a forbidden error
func ForbiddenServiceError(message string) error {
	return New("ERR_FORBIDDEN", message, 403)
}

// ConflictError creates a conflict error
func ConflictServiceError(message string) error {
	return New("ERR_CONFLICT", message, 409)
}

// InternalError creates an internal server error
func InternalServiceError(message string) error {
	return New("ERR_INTERNAL", message, 500)
}

// DatabaseError creates a database error
func DatabaseError(err error) error {
	return NewWithError("ERR_DATABASE", "database operation failed", 500, err)
}

// -- Helper Functions --

// AsCustomError extracts CustomError from error chain
func AsCustomError(err error) (*CustomError, bool) {
	ce, ok := err.(*CustomError)
	return ce, ok
}

// IsCustomError checks if error is CustomError
func IsCustomError(err error) bool {
	_, ok := err.(*CustomError)
	return ok
}

// GetStatus returns HTTP status code from error
func GetStatus(err error) int {
	if ce, ok := AsCustomError(err); ok {
		return ce.Status
	}
	return 500 // default to internal error
}

// GetCode returns error code
func GetCode(err error) string {
	if ce, ok := AsCustomError(err); ok {
		return ce.Code
	}
	return "ERR_UNKNOWN"
}
