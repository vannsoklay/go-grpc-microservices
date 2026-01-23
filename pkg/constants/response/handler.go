package response

import (
	"hpkg/grpc"
	"hpkg/i18n"

	"github.com/gofiber/fiber/v3"
)

// Error sends a standardized error response with translation support
func Error(c fiber.Ctx, status int, code string, msg ...string) error {
	lang, _ := c.Locals("lang").(string)
	if lang == "" {
		lang = "en"
	}

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

func FromGRPC[T any](c fiber.Ctx, err error, data ...T) error {
	if err == nil {
		var resp T
		if len(data) > 0 {
			resp = data[0]
		}
		return Success(c, fiber.StatusOK, resp)
	}

	httpErr := grpc.ToGRPC(err)
	return Error(c, httpErr.Status, httpErr.Code)
}
