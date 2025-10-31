package repository

import (
	"context"
	"errors"

	"github.com/wildanasyrof/backend-topup/internal/domain/dto"
	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	apperror "github.com/wildanasyrof/backend-topup/pkg/apperr"
	"github.com/wildanasyrof/backend-topup/pkg/pagination"
	"gorm.io/gorm"
)

type PaymentMethodsRepository interface {
	FindAll(ctx context.Context, q dto.PaymentMethodListQuery) (items []*entity.PaymentMethod, meta pagination.Meta, err error)
	Create(ctx context.Context, paymentMethod *entity.PaymentMethod) error
	Update(ctx context.Context, paymentMethod *entity.PaymentMethod) error
	Delete(ctx context.Context, id uint64) error
	FindByID(ctx context.Context, id uint64) (*entity.PaymentMethod, error)
}

type paymentMethodsRepository struct {
	db *gorm.DB
}

func NewPaymentMethodsRepository(db *gorm.DB) PaymentMethodsRepository {
	return &paymentMethodsRepository{db: db}
}

func (p *paymentMethodsRepository) Create(ctx context.Context, paymentMethod *entity.PaymentMethod) error {
	err := p.db.WithContext(ctx).Create(paymentMethod).Error

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return apperror.ErrConflict
	}

	return err
}

func (p *paymentMethodsRepository) Delete(ctx context.Context, id uint64) error {
	return p.db.WithContext(ctx).Delete(&entity.PaymentMethod{}, id).Error
}

func (p *paymentMethodsRepository) FindAll(ctx context.Context, q dto.PaymentMethodListQuery) (items []*entity.PaymentMethod, meta pagination.Meta, err error) {
	q.Normalize() // Terapkan DefaultPage dan DefaultLimit

	base := p.db.WithContext(ctx).
		Model(&entity.PaymentMethod{}).
		Preload("Provider") // <-- Ini sudah ada, bagus!

	// Tentukan kolom yang boleh di-sort
	allowedSort := map[string]struct{}{"created_at": {}, "name": {}, "type": {}, "id": {}}

	// Terapkan filter dan search
	filtered := base.
		Scopes(
			PaymentMethodFilters(q), // Filter by type, name, provider_id
			ILike([]string{"payment_methods.name", "payment_methods.type"}, q.Q), // Search by name dan type
		)

	// Hitung total data
	var total int64
	if err = filtered.Count(&total).Error; err != nil {
		return
	}

	// Ambil data untuk halaman ini
	if err = filtered.
		Scopes(
			func(db *gorm.DB) *gorm.DB { return pagination.ScopeSort(db, q.Sort, allowedSort) },
			func(db *gorm.DB) *gorm.DB { return pagination.ScopePaginate(db, q.Page, q.Limit) },
		).
		Find(&items).Error; err != nil {
		return
	}

	// Hitung metadata paginasi
	meta = pagination.CalcMeta(int(total), q.Page, q.Limit)
	return
}

func (p *paymentMethodsRepository) FindByID(ctx context.Context, id uint64) (*entity.PaymentMethod, error) {
	var data entity.PaymentMethod
	err := p.db.WithContext(ctx).First(&data, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperror.ErrNotFound
	}
	return &data, err
}

func (p *paymentMethodsRepository) Update(ctx context.Context, paymentMethod *entity.PaymentMethod) error {
	return p.db.WithContext(ctx).Save(paymentMethod).Error
}

func PaymentMethodFilters(q dto.PaymentMethodListQuery) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if q.ProviderID != nil {
			db = db.Where("provider_id = ?", *q.ProviderID)
		}
		if q.Type != nil {
			// Gunakan ILIKE agar lebih fleksibel
			db = db.Where("type ILIKE ?", "%"+*q.Type+"%")
		}
		if q.Name != nil {
			// Gunakan ILIKE agar lebih fleksibel
			db = db.Where("name ILIKE ?", "%"+*q.Name+"%")
		}
		return db
	}
}
