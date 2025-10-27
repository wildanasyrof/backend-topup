// internal/http/server/error_handler.go
package server

import (
	"github.com/gofiber/fiber/v2"
	apperror "github.com/wildanasyrof/backend-topup/pkg/apperr"
	"github.com/wildanasyrof/backend-topup/pkg/logger"
	"github.com/wildanasyrof/backend-topup/pkg/response"
)

func ErrorHandler(log logger.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		mapping, ae := apperror.MapToHTTP(err)

		// log: include request id, method, path, status, code
		log = log.With(map[string]any{
			"request_id": c.Get("X-Request-ID"),
			"method":     c.Method(),
			"path":       c.Path(),
			"status":     mapping.Status,
		})
		// internal/http/server/error_handler.go

		if ae != nil {
			log = log.With(map[string]any{"code": ae.Code})
			if ae.Cause != nil {
				// use Error(err, msg)
				log.Error(ae.Cause, "request failed")
			} else {
				log.Warn("request failed: " + ae.Msg)
			}
		} else {
			log.Error(err, "request failed")
		}

		env := response.Envelope{
			Success:   false,
			RequestID: c.Get("X-Request-ID"),
			Error: &response.ErrBody{
				Code:    string(aeCode(ae)),
				Message: publicMessage(mapping, ae),
				Fields:  fields(ae),
			},
		}
		return c.Status(mapping.Status).JSON(env)
	}
}

func aeCode(ae *apperror.AppError) apperror.Code {
	if ae == nil {
		return apperror.CodeInternal
	}
	return ae.Code
}
func publicMessage(m apperror.HTTPMapping, ae *apperror.AppError) string {
	if ae != nil && ae.Msg != "" {
		return ae.Msg
	}
	return m.Public
}
func fields(ae *apperror.AppError) map[string]string {
	if ae != nil && len(ae.Fields) > 0 {
		return ae.Fields
	}
	return nil
}
