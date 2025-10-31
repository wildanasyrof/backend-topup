package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wildanasyrof/backend-topup/internal/service"
	apperror "github.com/wildanasyrof/backend-topup/pkg/apperr"
	"github.com/wildanasyrof/backend-topup/pkg/response"
)

type SessionHandler struct {
	sessionSvc service.SessionService
}

func NewSessionHandler(sessionSvc service.SessionService) *SessionHandler {
	return &SessionHandler{sessionSvc: sessionSvc}
}

// GetMySessions mengambil daftar sesi aktif milik user
func (h *SessionHandler) GetMySessions(c *fiber.Ctx) error {
	uid, ok := c.Locals("user_id").(uint64)
	if !ok {
		return apperror.ErrUnauthorized
	}

	sessions, err := h.sessionSvc.GetByUserID(c.UserContext(), uid)
	if err != nil {
		return err
	}

	return response.OK(c, sessions)
}

// RevokeSession meng-invalidate sebuah refresh token (sesi)
func (h *SessionHandler) RevokeSession(c *fiber.Ctx) error {
	uid, ok := c.Locals("user_id").(uint64)
	if !ok {
		return apperror.ErrUnauthorized
	}

	sessionID := c.Params("id")
	if sessionID == "" {
		return apperror.New(apperror.CodeBadRequest, "session id is required", nil)
	}

	// Panggil service untuk me-revoke
	if err := h.sessionSvc.RevokeSession(c.UserContext(), uid, sessionID); err != nil {
		return err
	}

	return response.OK(c, fiber.Map{"message": "session revoked"})
}
