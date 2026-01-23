package response

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message,omitempty"` // optional message, useful for dev/debug
}

// ResponseEnvelope wraps all API responses
type ResponseEnvelope[T any] struct {
	Success bool           `json:"success"`
	Data    T              `json:"data,omitempty"`
	Error   *ErrorResponse `json:"error,omitempty"`
}
