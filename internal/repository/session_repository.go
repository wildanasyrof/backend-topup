package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	apperror "github.com/wildanasyrof/backend-topup/pkg/apperr"
	"gorm.io/gorm"
)

type SessionRepository interface {
	Create(ctx context.Context, session *entity.UserSession) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.UserSession, error)
	Delete(ctx context.Context, id uuid.UUID) error
	FindByUserID(ctx context.Context, userID uint64) ([]*entity.UserSession, error)
	Revoke(ctx context.Context, id uuid.UUID) error
}

type sessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) Create(ctx context.Context, session *entity.UserSession) error {
	return r.db.WithContext(ctx).Create(session).Error
}

func (r *sessionRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.UserSession, error) {
	var session entity.UserSession
	// Kita Preload User karena kita perlu Role untuk membuat Access Token baru
	err := r.db.WithContext(ctx).Preload("User").Where("id = ?", id).First(&session).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.ErrUnauthorized // Token tidak ditemukan = tidak valid
		}
		return nil, err
	}
	return &session, nil
}

func (r *sessionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.UserSession{}, id).Error
}

func (r *sessionRepository) FindByUserID(ctx context.Context, userID uint64) ([]*entity.UserSession, error) {
	var sessions []*entity.UserSession
	// Hanya ambil sesi yang belum di-revoke dan belum kedaluwarsa
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_revoked = ? AND expires_at > ?", userID, false, time.Now()).
		Order("created_at desc").
		Find(&sessions).Error
	return sessions, err
}

func (r *sessionRepository) Revoke(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&entity.UserSession{}).Where("id = ?", id).
		Update("is_revoked", true).Error
}
