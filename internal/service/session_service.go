package service

import (
	"context"
	"errors" // <-- Import

	"github.com/google/uuid"
	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	"github.com/wildanasyrof/backend-topup/internal/repository"
	apperror "github.com/wildanasyrof/backend-topup/pkg/apperr"
	"gorm.io/gorm" // <-- Import
)

type SessionService interface {
	GetByUserID(ctx context.Context, userID uint64) ([]*entity.UserSession, error)
	RevokeSession(ctx context.Context, userID uint64, sessionToRevokeID string) error
}

type sessionService struct {
	sessionRepo repository.SessionRepository
}

func NewSessionService(sessionRepo repository.SessionRepository) SessionService {
	return &sessionService{sessionRepo: sessionRepo}
}

func (s *sessionService) GetByUserID(ctx context.Context, userID uint64) ([]*entity.UserSession, error) {
	return s.sessionRepo.FindByUserID(ctx, userID)
}

func (s *sessionService) RevokeSession(ctx context.Context, userID uint64, sessionToRevokeID string) error {
	sessionID, err := uuid.Parse(sessionToRevokeID)
	if err != nil {
		return apperror.New(apperror.CodeBadRequest, "invalid session id format", err)
	}

	// Verifikasi kepemilikan
	session, err := s.sessionRepo.FindByID(ctx, sessionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.ErrNotFound
		}
		return err
	}

	if session.UserID != userID {
		// User mencoba menghapus sesi orang lain
		return apperror.ErrForbidden
	}

	// Kita tandai sebagai revoked, bukan dihapus
	return s.sessionRepo.Revoke(ctx, sessionID)
}
