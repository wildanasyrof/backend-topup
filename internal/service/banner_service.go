package service

import (
	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	"github.com/wildanasyrof/backend-topup/internal/repository"
)

type BannerService interface {
	Create(imgUrl string) (*entity.Banner, error)
	FindAll() ([]*entity.Banner, error)
	Update(id int, imgUrl string) (*entity.Banner, error)
	Delete(id int) (*entity.Banner, error)
}

type bannerService struct {
	bannerRepo repository.BannerRepository
}

func NewBannerService(bannerRepo repository.BannerRepository) BannerService {
	return &bannerService{
		bannerRepo: bannerRepo,
	}
}

// Create implements BannerService.
func (b *bannerService) Create(imgUrl string) (*entity.Banner, error) {
	banner := &entity.Banner{
		ImgUrl: imgUrl,
	}

	if err := b.bannerRepo.Create(banner); err != nil {
		return nil, err
	}

	return banner, nil
}

// Delete implements BannerService.
func (b *bannerService) Delete(id int) (*entity.Banner, error) {
	banner, err := b.bannerRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if err := b.bannerRepo.Delete(id); err != nil {
		return nil, err
	}

	return banner, nil
}

// FindAll implements BannerService.
func (b *bannerService) FindAll() ([]*entity.Banner, error) {
	banners, err := b.bannerRepo.FindAll()
	if err != nil {
		return nil, err
	}

	return banners, nil
}

// Update implements BannerService.
func (b *bannerService) Update(id int, imgUrl string) (*entity.Banner, error) {
	banner, err := b.bannerRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	banner.ImgUrl = imgUrl

	if err := b.bannerRepo.Update(banner); err != nil {
		return nil, err
	}

	return banner, nil
}
