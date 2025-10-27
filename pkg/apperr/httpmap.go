package apperror

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

type HTTPMapping struct {
	Status int
	Public string // message exposed to clients
}

func MapToHTTP(err error) (HTTPMapping, *AppError) {
	// default
	mapping := HTTPMapping{Status: fiber.StatusInternalServerError, Public: "Something went wrong"}

	var ae *AppError
	if errors.As(err, &ae) {
		switch ae.Code {
		case CodeBadRequest, CodeUnprocessable:
			mapping.Status = fiber.StatusBadRequest
			mapping.Public = ae.Msg
		case CodeUnauthorized:
			mapping.Status = fiber.StatusUnauthorized
			mapping.Public = "Unauthorized"
		case CodeForbidden:
			mapping.Status = fiber.StatusForbidden
			mapping.Public = "Forbidden"
		case CodeNotFound:
			mapping.Status = fiber.StatusNotFound
			mapping.Public = ae.Msg
		case CodeConflict:
			mapping.Status = fiber.StatusConflict
			mapping.Public = ae.Msg
		case CodeRateLimited:
			mapping.Status = fiber.StatusTooManyRequests
			mapping.Public = "Too many requests"
		case CodeUnavailable:
			mapping.Status = fiber.StatusServiceUnavailable
			mapping.Public = "Service unavailable"
		case CodeTimeout:
			mapping.Status = fiber.StatusGatewayTimeout
			mapping.Public = "Request timeout"
		case CodeInternal:
			fallthrough
		default:
			mapping.Status = fiber.StatusInternalServerError
			mapping.Public = "Internal server error"
		}
		return mapping, ae
	}

	// Unknown non-AppError
	return mapping, nil
}
