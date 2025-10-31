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

type ProviderRepository interface {
	Create(ctx context.Context, req *entity.Provider) error
	FindAll(ctx context.Context, q dto.ProviderListQuery) (items []*entity.Provider, meta pagination.Meta, err error)
	FindByID(ctx context.Context, id int64) (*entity.Provider, error)
	FindBySlug(ctx context.Context, slug string) (*entity.Provider, error)
	Update(ctx context.Context, req *entity.Provider) error
	Delete(ctx context.Context, id int64) error
}

type providerRepository struct {
	db *gorm.DB
}

func NewProviderRepository(db *gorm.DB) ProviderRepository {
	return &providerRepository{db: db}
}

// Create implements ProviderRepository.
func (p *providerRepository) Create(ctx context.Context, req *entity.Provider) error {
	err := p.db.WithContext(ctx).Create(req).Error

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return apperror.ErrConflict
	}

	return err
}

// Delete implements ProviderRepository.
func (p *providerRepository) Delete(ctx context.Context, id int64) error {
	return p.db.WithContext(ctx).Delete(&entity.Provider{}, id).Error
}

// FindAll implements ProviderRepository.
func (p *providerRepository) FindAll(ctx context.Context, q dto.ProviderListQuery) (items []*entity.Provider, meta pagination.Meta, err error) {
	q.Normalize() // Terapkan DefaultPage dan DefaultLimit

	base := p.db.WithContext(ctx).Model(&entity.Provider{})

	// Tentukan kolom yang boleh di-sort
	allowedSort := map[string]struct{}{"created_at": {}, "name": {}, "ref": {}, "id": {}}

	// Terapkan search
	filtered := base.
		Scopes(
			ILike([]string{"providers.name", "providers.ref"}, q.Q), // Search by name dan ref
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

// FindByID implements ProviderRepository.
func (p *providerRepository) FindByID(ctx context.Context, id int64) (*entity.Provider, error) {
	var provider entity.Provider
	err := p.db.WithContext(ctx).Where("id = ?", id).First(&provider).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperror.ErrNotFound
	}

	return &provider, err
}

func (p *providerRepository) FindBySlug(ctx context.Context, slug string) (*entity.Provider, error) {
	var provider entity.Provider
	err := p.db.WithContext(ctx).Where("slug = ?", slug).First(&provider).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperror.ErrNotFound
	}

	return &provider, err
}

// Update implements ProviderRepository.
func (p *providerRepository) Update(ctx context.Context, req *entity.Provider) error {
	return p.db.WithContext(ctx).Save(req).Error
}
