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

type ProductRepository interface {
	Create(ctx context.Context, req *entity.Product) error
	FindAll(ctx context.Context, q dto.ProductListQuery) (items []entity.Product, meta pagination.Meta, err error)
	FindByID(ctx context.Context, id int) (*entity.Product, error)
	FindByCode(ctx context.Context, code string) (*entity.Product, error)
	Update(ctx context.Context, req *entity.Product) error
	Delete(ctx context.Context, id int) error
	UpsertProductByCode(ctx context.Context, req *entity.Product) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

// Create implements ProductRepository.
func (p *productRepository) Create(ctx context.Context, req *entity.Product) error {
	err := p.db.WithContext(ctx).Create(req).Error

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return apperror.New(apperror.CodeConflict, "product already exist", err)
	}

	return err
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
	// TAMBAHKAN kolom baru ke allowedSort
	allowedSort := map[string]struct{}{"created_at": {}, "name": {}, "price": {}, "id": {}, "stock": {}, "seller_name": {}}

	filtered := base.
		Scopes(
			// --- UBAH DI SINI ---
			// Teruskan seluruh DTO query ke ProductFilters
			ProductFilters(q),
			// Perluas pencarian 'q' untuk mencakup sku_code dan seller_name
			ILike([]string{"products.name", "products.sku_code", "products.seller_name"}, q.Q),
			// ---
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

	err := p.db.WithContext(ctx).Preload("Prices").First(&product, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperror.ErrNotFound
	}

	return &product, nil
}

// Update implements ProductRepository.
func (p *productRepository) Update(ctx context.Context, req *entity.Product) error {
	return p.db.WithContext(ctx).Save(req).Error
}

// Resource filter
func ProductFilters(q dto.ProductListQuery) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if q.ProviderID != nil {
			db = db.Where("products.provider_id = ?", *q.ProviderID)
		}
		if q.CategoryID != nil {
			db = db.Where("products.category_id = ?", *q.CategoryID)
		}
		if q.Status != nil {
			db = db.Where("products.status = ?", *q.Status)
		}
		if q.SellerName != nil {
			// Gunakan ILIKE untuk seller_name agar lebih fleksibel
			db = db.Where("products.seller_name ILIKE ?", "%"+*q.SellerName+"%")
		}
		if q.SkuCode != nil {
			db = db.Where("products.sku_code = ?", *q.SkuCode)
		}

		if q.LevelID != nil {
			// example: join prices table to ensure at least one price for level
			db = db.Joins("LEFT JOIN prices ON prices.product_id = products.id AND prices.user_level_id = ?", *q.LevelID).
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

// FindByCode implements ProductRepository.
func (p *productRepository) FindByCode(ctx context.Context, code string) (*entity.Product, error) {
	var product entity.Product
	err := p.db.WithContext(ctx).Where("sku_code = ?", code).First(&product).Error

	return &product, err
}

// UpsertProductByCode implements ProductRepository.
func (p *productRepository) UpsertProductByCode(ctx context.Context, req *entity.Product) error {
	var isExist entity.Product

	err := p.db.WithContext(ctx).Where("sku_code = ?", req.SkuCode).First(&isExist).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return p.db.WithContext(ctx).Create(req).Error
	}

	if err != nil {
		return err
	}

	return p.db.WithContext(ctx).Save(req).Error
}
