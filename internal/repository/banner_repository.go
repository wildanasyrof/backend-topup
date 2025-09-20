package repository

import (
	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	"gorm.io/gorm"
)

type BannerRepository interface {
	Create(req *entity.Banner) error
	FindAll() ([]*entity.Banner, error)
	FindByID(id int) (*entity.Banner, error)
	Update(req *entity.Banner) error
	Delete(id int) error
}

type bannerRepository struct {
	db *gorm.DB
}

func NewBannerRepository(db *gorm.DB) BannerRepository {
	return &bannerRepository{db: db}
}

// Create implements BannerRepository.
func (b *bannerRepository) Create(req *entity.Banner) error {
	return b.db.Create(req).Error
}

// Delete implements BannerRepository.
func (b *bannerRepository) Delete(id int) error {
	return b.db.Delete(&entity.Banner{}, id).Error
}

// FindAll implements BannerRepository.
func (b *bannerRepository) FindAll() ([]*entity.Banner, error) {
	var banners []*entity.Banner
	if err := b.db.Find(&banners).Error; err != nil {
		return nil, err
	}

	return banners, nil
}

// FindByID implements BannerRepository.
func (b *bannerRepository) FindByID(id int) (*entity.Banner, error) {
	var banner entity.Banner
	if err := b.db.First(&banner, id).Error; err != nil {
		return nil, err
	}
	return &banner, nil
}

// Update implements BannerRepository.
func (b *bannerRepository) Update(req *entity.Banner) error {
	return b.db.Save(req).Error
}
