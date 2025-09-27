package repository

import (
	"context"

	"github.com/wildanasyrof/backend-topup/internal/domain/dto"
	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	"github.com/wildanasyrof/backend-topup/pkg/pagination"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(ctx context.Context, req *entity.Product) error
	FindAll(ctx context.Context, q dto.ProductListQuery) (items []entity.Product, meta pagination.Meta, err error)
	FindByID(ctx context.Context, id int) (*entity.Product, error)
	Update(ctx context.Context, req *entity.Product) error
	Delete(ctx context.Context, id int) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

// Create implements ProductRepository.
func (p *productRepository) Create(ctx context.Context, req *entity.Product) error {
	return p.db.WithContext(ctx).Create(req).Error
}

// Delete implements ProductRepository.
func (p *productRepository) Delete(ctx context.Context, id int) error {
	return p.db.WithContext(ctx).Delete(&entity.Product{}, id).Error
}

// FindAll implements ProductRepository.
func (p *productRepository) FindAll(ctx context.Context, q dto.ProductListQuery) (items []entity.Product, meta pagination.Meta, err error) {
	q.Normalize()

	base := p.db.WithContext(ctx).
		Model(&entity.Product{}).
		Preload("Prices").Preload("Provider").Preload("Category")

	// count with filters/search
	allowedSort := map[string]struct{}{"created_at": {}, "name": {}, "price": {}, "id": {}}
	filtered := base.
		Scopes(
			ProductFilters(q.ProviderID, q.CategoryID, q.LevelID, q.Active),
			ILike([]string{"products.name"}, q.Q),
		)

	var total int64
	if err = filtered.Count(&total).Error; err != nil {
		return
	}

	// data page
	if err = filtered.
		Scopes(
			func(db *gorm.DB) *gorm.DB { return pagination.ScopeSort(db, q.Sort, allowedSort) },
			func(db *gorm.DB) *gorm.DB { return pagination.ScopePaginate(db, q.Page, q.Limit) },
		).
		Find(&items).Error; err != nil {
		return
	}

	meta = pagination.CalcMeta(int(total), q.Page, q.Limit)
	return
}

func (p *productRepository) FindByID(ctx context.Context, id int) (*entity.Product, error) {
	var product entity.Product

	if err := p.db.WithContext(ctx).Preload("Prices").First(&product, id).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

// Update implements ProductRepository.
func (p *productRepository) Update(ctx context.Context, req *entity.Product) error {
	return p.db.WithContext(ctx).Save(req).Error
}

// Resource filter

func ProductFilters(providerID, categoryID, levelID *uint, active *bool) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if providerID != nil {
			db = db.Where("provider_id = ?", *providerID)
		}
		if categoryID != nil {
			db = db.Where("category_id = ?", *categoryID)
		}
		if active != nil {
			db = db.Where("active = ?", *active)
		}
		if levelID != nil {
			// example: join prices table to ensure at least one price for level
			db = db.Joins("LEFT JOIN prices ON prices.product_id = products.id AND prices.user_level_id = ?", *levelID).
				Where("prices.id IS NOT NULL")
		}
		return db
	}
}

// Search query

func ILike(cols []string, q string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if q == "" {
			return db
		}
		like := "%" + q + "%"
		// (name ILIKE ? OR code ILIKE ? ...)
		cond := db
		for i, col := range cols {
			if i == 0 {
				cond = db.Where(col+" ILIKE ?", like)
			} else {
				cond = cond.Or(col+" ILIKE ?", like)
			}
		}
		return cond
	}
}
