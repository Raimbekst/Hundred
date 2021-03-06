package service

import (
	"HundredToFive/internal/domain"
	"HundredToFive/internal/repository"
	"HundredToFive/pkg/media"
	"fmt"
)

type BannerService struct {
	repos repository.Banner
}

func NewBannerService(repos repository.Banner) *BannerService {
	return &BannerService{repos: repos}
}

func (b *BannerService) Create(banner domain.Banner) (int, error) {
	return b.repos.Create(banner)
}
func (b *BannerService) GetAll(page domain.Pagination, status int, lang string) (*domain.GetAllBannersCategoryResponse, error) {
	return b.repos.GetAll(page, status, lang)
}
func (b *BannerService) GetById(id int) (domain.Banner, error) {
	return b.repos.GetById(id)
}
func (b *BannerService) Update(id int, inp domain.Banner) error {
	image, err := b.repos.Update(id, inp)

	if err != nil {
		return fmt.Errorf("service.Update: %w", err)
	}
	if image != "" {
		err = media.DeleteImage(image)
		if err != nil {
			return fmt.Errorf("service.Update: %w", err)
		}
	}
	return nil
}
func (b *BannerService) Delete(id int) error {
	img, err := b.repos.Delete(id)
	if err != nil {
		return fmt.Errorf("service.Delete: %w", err)

	}
	if img != "" {
		err = media.DeleteImage(img)
		if err != nil {
			return fmt.Errorf("service.Delete: %w", err)
		}
	}
	return nil
}
