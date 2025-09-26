package repository

import (
	"context"

	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	"gorm.io/gorm"
)

type BannerRepository interface {
	Create(ctx context.Context, req *entity.Banner) error
	FindAll(ctx context.Context) ([]*entity.Banner, error)
	FindByID(ctx context.Context, id int) (*entity.Banner, error)
	Update(ctx context.Context, req *entity.Banner) error
	Delete(ctx context.Context, id int) error
}

type bannerRepository struct {
	db *gorm.DB
}

func NewBannerRepository(db *gorm.DB) BannerRepository {
	return &bannerRepository{db: db}
}

// Create implements BannerRepository.
func (b *bannerRepository) Create(ctx context.Context, req *entity.Banner) error {
	return b.db.WithContext(ctx).Create(req).Error
}

// Delete implements BannerRepository.
func (b *bannerRepository) Delete(ctx context.Context, id int) error {
	return b.db.WithContext(ctx).Delete(&entity.Banner{}, id).Error
}

// FindAll implements BannerRepository.
func (b *bannerRepository) FindAll(ctx context.Context) ([]*entity.Banner, error) {
	var banners []*entity.Banner
	if err := b.db.WithContext(ctx).Find(&banners).Error; err != nil {
		return nil, err
	}

	return banners, nil
}

// FindByID implements BannerRepository.
func (b *bannerRepository) FindByID(ctx context.Context, id int) (*entity.Banner, error) {
	var banner entity.Banner
	if err := b.db.WithContext(ctx).First(&banner, id).Error; err != nil {
		return nil, err
	}
	return &banner, nil
}

// Update implements BannerRepository.
func (b *bannerRepository) Update(ctx context.Context, req *entity.Banner) error {
	return b.db.WithContext(ctx).Save(req).Error
}
