package service

import (
	"context"

	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	"github.com/wildanasyrof/backend-topup/internal/repository"
)

type BannerService interface {
	Create(ctx context.Context, imgUrl string) (*entity.Banner, error)
	FindAll(ctx context.Context) ([]*entity.Banner, error)
	Update(ctx context.Context, id int, imgUrl string) (*entity.Banner, error)
	Delete(ctx context.Context, id int) (*entity.Banner, error)
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
func (b *bannerService) Create(ctx context.Context, imgUrl string) (*entity.Banner, error) {
	banner := &entity.Banner{
		ImgUrl: imgUrl,
	}

	if err := b.bannerRepo.Create(ctx, banner); err != nil {
		return nil, err
	}

	return banner, nil
}

// Delete implements BannerService.
func (b *bannerService) Delete(ctx context.Context, id int) (*entity.Banner, error) {
	banner, err := b.bannerRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := b.bannerRepo.Delete(ctx, id); err != nil {
		return nil, err
	}

	return banner, nil
}

// FindAll implements BannerService.
func (b *bannerService) FindAll(ctx context.Context) ([]*entity.Banner, error) {
	banners, err := b.bannerRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return banners, nil
}

// Update implements BannerService.
func (b *bannerService) Update(ctx context.Context, id int, imgUrl string) (*entity.Banner, error) {
	banner, err := b.bannerRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	banner.ImgUrl = imgUrl

	if err := b.bannerRepo.Update(ctx, banner); err != nil {
		return nil, err
	}

	return banner, nil
}
