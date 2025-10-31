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

type PriceRepository interface {
	Create(ctx context.Context, req *entity.Price) error
	FindAll(ctx context.Context, q dto.PriceListQuery) (items []*entity.Price, meta pagination.Meta, err error)
	FindByID(ctx context.Context, id int) (*entity.Price, error)
	FindByProductIDnUserLevelID(ctx context.Context, productId int, userLevelId int) (*entity.Price, error)
	Update(ctx context.Context, price *entity.Price) error
	Delete(ctx context.Context, id int) error
}

type priceRepository struct {
	db *gorm.DB
}

func NewPriceRepository(db *gorm.DB) PriceRepository {
	return &priceRepository{db: db}
}

// Create implements PriceRepository.
func (p *priceRepository) Create(ctx context.Context, req *entity.Price) error {
	return p.db.WithContext(ctx).Create(req).Error
}

// Delete implements PriceRepository.
func (p *priceRepository) Delete(ctx context.Context, id int) error {
	return p.db.WithContext(ctx).Delete(&entity.Price{}, id).Error
}

// FindAll implements PriceRepository.
func (p *priceRepository) FindAll(ctx context.Context, q dto.PriceListQuery) (items []*entity.Price, meta pagination.Meta, err error) {
	q.Normalize() // Terapkan DefaultPage dan DefaultLimit

	base := p.db.WithContext(ctx).
		Model(&entity.Price{}).
		Preload("UserLevel"). // <-- PENTING: Ambil data UserLevel
		Preload("Product")    // <-- PENTING: Ambil data Product

	// Tentukan kolom yang boleh di-sort
	allowedSort := map[string]struct{}{"created_at": {}, "amount": {}, "id": {}}

	// Terapkan filter
	filtered := base.
		Scopes(
			PriceFilters(q), // Terapkan filter by product_id / user_level_id
		)
	// Tidak perlu ILike (search 'q') karena tidak ada kolom teks yang relevan

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

// FindByID implements PriceRepository.
func (p *priceRepository) FindByID(ctx context.Context, id int) (*entity.Price, error) {
	var price entity.Price

	err := p.db.WithContext(ctx).Where("id = ?", id).First(&price).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperror.ErrNotFound
	}

	return &price, err
}

// Update implements PriceRepository.
func (p *priceRepository) Update(ctx context.Context, price *entity.Price) error {
	return p.db.WithContext(ctx).Save(price).Error
}

func (p *priceRepository) FindByProductIDnUserLevelID(ctx context.Context, productId int, userLevelId int) (*entity.Price, error) {
	var price entity.Price

	err := p.db.WithContext(ctx).
		Where("product_id = ? AND user_level_id = ?", productId, userLevelId).
		First(&price).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperror.ErrNotFound
	}

	return &price, err
}

func PriceFilters(q dto.PriceListQuery) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if q.ProductID != nil {
			db = db.Where("product_id = ?", *q.ProductID)
		}
		if q.UserLevelID != nil {
			db = db.Where("user_level_id = ?", *q.UserLevelID)
		}
		return db
	}
}
