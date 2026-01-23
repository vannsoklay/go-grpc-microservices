package constants

// ===== Auth Errors =====
const (
	ErrAuthHeaderMissingCode = "AUTH_HEADER_MISSING"
	ErrAuthHeaderMissingMsg  = "Authorization header is required"

	ErrTokenInvalidCode = "TOKEN_INVALID"
	ErrTokenInvalidMsg  = "Invalid or malformed authorization token"

	ErrTokenExpiredCode = "TOKEN_EXPIRED"
	ErrTokenExpiredMsg  = "Your authentication token has expired. Please log in again"

	ErrUnauthorizedCode = "UNAUTHORIZED"
	ErrUnauthorizedMsg  = "Unauthorized request. Please provide valid credentials"

	ErrForbiddenCode = "FORBIDDEN"
	ErrForbiddenMsg  = "You do not have permission to access this resource"

	ErrNotFoundCode = "NOT_FOUND"
	ErrNotFoundMsg  = "The requested resource was not found"

	ErrServiceUnavailableCode = "SERVICE_UNAVAILABLE"
	ErrServiceUnavailableMsg  = "The service is temporarily unavailable. Please try again later"
)

// ===== Shop Errors =====
const (
	ShopRequiredCode = "SHOP_REQUIRED"
	ShopRequiredMsg  = "X-Shop-Id header is required"

	ShopAccessDeniedCode = "SHOP_ACCESS_DENIED"
	ShopAccessDeniedMsg  = "You do not have access to this shop"

	ShopNotFoundCode = "SHOP_NOT_FOUND"
	ShopNotFoundMsg  = "The requested shop was not found"
)

// ===== Validation Errors =====
const (
	ErrInvalidInputCode = "INVALID_INPUT"
	ErrInvalidInputMsg  = "Request contains invalid input data"

	ErrBadRequestCode = "BAD_REQUEST"
	ErrBadRequestMsg  = "The request is malformed or missing required fields"

	ErrInvalidPayloadCode = "INVALID_PAYLOAD"
	ErrInvalidPayloadMsg  = "Request payload is invalid or could not be parsed"
)

// ===== Service Errors =====
const (
	ErrUserServiceCode = "USER_SERVICE_ERROR"
	ErrUserServiceMsg  = "Failed to fetch user data. Please try again later"

	ErrDatabaseCode = "DATABASE_ERROR"
	ErrDatabaseMsg  = "A database error occurred. Please try again later"

	ErrInternalCode = "INTERNAL_SERVER_ERROR"
	ErrInternalMsg  = "An unexpected internal server error occurred. Our team has been notified"
)

// ===== Success Responses =====
const (
	SuccessCode = ""
	SuccessMsg  = "Success"
)
