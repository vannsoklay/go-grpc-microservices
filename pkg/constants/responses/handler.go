package responses

import (
	"errors"
	"hpkg/i18n"

	"github.com/gofiber/fiber/v3"
)

// Error sends a standardized error response with translation support
func Error(c fiber.Ctx, status int, code string, msg ...string) error {
	lang := getLanguage(c)

	var message string
	if len(msg) > 0 && msg[0] != "" {
		message = msg[0] // custom message passed
	} else {
		message = i18n.Translate(code, lang)
	}

	return c.Status(status).JSON(ResponseEnvelope[any]{
		Success: false,
		Error: &ErrorResponse{
			Code:    code,
			Message: message,
		},
	})
}

// Success sends a standardized success response
func Success[T any](c fiber.Ctx, status int, data T) error {
	return c.Status(status).JSON(ResponseEnvelope[T]{
		Success: true,
		Data:    data,
	})
}

// Created sends a 201 Created response
func Created[T any](c fiber.Ctx, data T) error {
	return Success(c, fiber.StatusCreated, data)
}

// OK sends a 200 OK response
func OK[T any](c fiber.Ctx, data T) error {
	return Success(c, fiber.StatusOK, data)
}

// NoContent sends a 204 No Content response
func NoContent(c fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

// BadRequest sends a 400 Bad Request error
func BadRequest(c fiber.Ctx, code string, msg ...string) error {
	return Error(c, fiber.StatusBadRequest, code, msg...)
}

// Unauthorized sends a 401 Unauthorized error
func Unauthorized(c fiber.Ctx, code string, msg ...string) error {
	return Error(c, fiber.StatusUnauthorized, code, msg...)
}

// Forbidden sends a 403 Forbidden error
func Forbidden(c fiber.Ctx, code string, msg ...string) error {
	return Error(c, fiber.StatusForbidden, code, msg...)
}

// NotFound sends a 404 Not Found error
func NotFound(c fiber.Ctx, code string, msg ...string) error {
	return Error(c, fiber.StatusNotFound, code, msg...)
}

// Conflict sends a 409 Conflict error
func Conflict(c fiber.Ctx, code string, msg ...string) error {
	return Error(c, fiber.StatusConflict, code, msg...)
}

// InternalError sends a 500 Internal Server Error
func InternalError(c fiber.Ctx, code string, msg ...string) error {
	return Error(c, fiber.StatusInternalServerError, code, msg...)
}

// FromGRPC converts gRPC error to HTTP response
func FromGRPC[T any](c fiber.Ctx, err error, data ...T) error {
	if err == nil {
		var resp T
		if len(data) > 0 {
			resp = data[0]
		}
		return OK(c, resp)
	}

	httpErr := ToGRPC(err)
	return Error(c, httpErr.Status, httpErr.Code)
}

// FromError handles generic errors and converts to appropriate HTTP response
func FromError(c fiber.Ctx, err error) error {
	if err == nil {
		return NoContent(c)
	}

	// Handle specific error types
	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		return Error(c, fiberErr.Code, "error", fiberErr.Message)
	}

	// Try gRPC conversion
	httpErr := ToGRPC(err)
	return Error(c, httpErr.Status, httpErr.Code)
}

// Paginated sends a paginated response
type PaginatedData[T any] struct {
	Items      []T   `json:"items"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}

func Paginated[T any](c fiber.Ctx, items []T, total int64, page int, pageSize int) error {
	totalPages := (int(total) + pageSize - 1) / pageSize
	if totalPages < 1 {
		totalPages = 1
	}

	return OK(c, PaginatedData[T]{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

// getLanguage extracts language from Fiber context with fallback
func getLanguage(c fiber.Ctx) string {
	// Try from locals first
	if lang, ok := c.Locals("lang").(string); ok && lang != "" {
		return lang
	}

	// Try from query param
	if lang := c.Query("lang"); lang != "" {
		return lang
	}

	// Try from header
	if lang := c.Get("Accept-Language"); lang != "" {
		return lang
	}

	// Default to English
	return "en"
}
