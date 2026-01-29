package responses

// ===== Auth Errors =====
const (
	ErrAuthHeaderMissingCode = "AUTH_HEADER_MISSING"
	ErrAuthHeaderMissingMsg  = "Authorization header is required"

	ErrTokenInvalidCode = "TOKEN_INVALID"
	ErrTokenInvalidMsg  = "Invalid token"

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

	ErrInvalidCredentialsCode = "INVALID_CREDENTIALS"
	ErrInvalidCredentialsMsg  = "Invalid email or password"

	ErrUserCreateFailedCode = "USER_CREATE_FAILED"
	ErrUserCreateFailedMsg  = "Failed to create user"

	TokenGenerateFailedCode       = "TOKEN_GENERATE_FAILED"
	AccessTokenGenerateFailedMsg  = "Failed to generate access token"
	RefreshTokenGenerateFailedMsg = "Failed to generate refresh token"
)

// ===== user Errors =====

const (
	UserNotFoundCode         = "USER_NOT_FOUND"
	UserFetchFailedCode      = "USER_FETCH_FAILED"
	UsernameExistsCode       = "USERNAME_EXISTS"
	UsernameUpdateFailedCode = "USERNAME_UPDATE_FAILED"
	UsernameCheckFailedCode  = "USERNAME_CHECK_FAILED"
	InvalidRequestCode       = "INVALID_REQUEST"
	UnauthenticatedCode      = "UNAUTHENTICATED"
	RequestCanceledCode      = "REQUEST_CANCELED"

	UserNotFoundMsg         = "User not found"
	UserFetchFailedMsg      = "Failed to fetch user"
	UsernameExistsMsg       = "Username already exists"
	UsernameUpdateFailedMsg = "Failed to update username"
	UsernameCheckFailedMsg  = "Failed to check username"
	InvalidRequestMsg       = "Invalid request"
	UnauthenticatedMsg      = "User not authenticated"
	RequestCanceledMsg      = "Request canceled"
)

// ===== Shop Errors =====
const (
	ShopRequiredCode = "SHOP_REQUIRED"
	ShopRequiredMsg  = "X-Shop-Id header is required"

	ShopCreateFailedCode = "SHOP_CREATE_FAILED"
	ShopCreateFailedMsg  = "Unable to create shop at this time. Please try again later"

	ShopSlugExistsCode = "SHOP_SLUG_EXISTS"
	ShopSlugExistsMsg  = "Shop slug already exists"

	ShopAccessDeniedCode = "SHOP_ACCESS_DENIED"
	ShopAccessDeniedMsg  = "You do not have access to this shop"

	ShopLimitExceededCode  = "SHOP_LIMIT_EXCEEDED"
	ShopLimitExceededMsg   = "You can only create up to 2 shops"
	ShopValidateFailedCode = "SHOP_VALIDATE_FAILED"
	ShopValidateFailedMsg  = "Failed to retrieve shop"

	ShopNotFoundCode = "SHOP_NOT_FOUND"
	ShopNotFoundMsg  = "Shop not found"

	ShopListFailedCode = "SHOP_LIST_FAILED"
	ShopListFailedMsg  = "Failed to retrieve shops"

	ShopUpdateFailedCode = "SHOP_UPDATE_FAILED"
	ShopUpdateFailedMsg  = "Failed to update shop"

	ShopDeleteFailedCode = "SHOP_DELETE_FAILED"
	ShopDeleteFailedMsg  = "Failed to delete shop"
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

// Product Service Errors
const (
	ErrProductServiceCode = "PRODUCT_SERVICE_ERROR"
	ErrProductServiceMsg  = "Failed to fetch product data. Please try again later"

	ErrProductNotFoundCode = "PRODUCT_NOT_FOUND"
	ErrProductNotFoundMsg  = "Product not found"

	ErrProductInvalidCode = "PRODUCT_INVALID"
	ErrProductInvalidMsg  = "Invalid product data provided"

	ErrProductOutOfStockCode = "PRODUCT_OUT_OF_STOCK"
	ErrProductOutOfStockMsg  = "Product is out of stock"

	ErrProductConflictCode = "PRODUCT_CONFLICT"
	ErrProductConflictMsg  = "Product already exists"

	ErrProductCategoryNotFoundCode = "PRODUCT_CATEGORY_NOT_FOUND"
	ErrProductCategoryNotFoundMsg  = "Product category not found"
)

// ===== Success Responses =====
const (
	SuccessCode = ""
	SuccessMsg  = "Success"
)
