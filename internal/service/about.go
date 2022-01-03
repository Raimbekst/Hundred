package service

import (
	"HundredToFive/internal/domain"
	"HundredToFive/internal/repository"
)

type AboutService struct {
	repos repository.About
}

func NewAboutService(repos repository.About) *AboutService {
	return &AboutService{repos: repos}
}

func (a *AboutService) Create(about domain.AboutUs) error {
	return a.repos.Create(about)
}

func (a *AboutService) GetAll() ([]*domain.AboutUs, error) {
	return a.repos.GetAll()
}

func (a *AboutService) Update(about domain.AboutUs) error {
	return a.repos.Update(about)
}

func (a *AboutService) Delete() error {
	return a.repos.Delete()
}
