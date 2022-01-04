package service

import (
	"HundredToFive/internal/domain"
	"HundredToFive/internal/repository"
)

type WinnerService struct {
	repos repository.Winner
}

func NewWinnerService(repos repository.Winner) *WinnerService {
	return &WinnerService{repos: repos}
}

func (w *WinnerService) CreateWinner(input domain.WinnerInput) error {
	return w.repos.CreateWinner(input)
}

func (w *WinnerService) GetAll(page domain.Pagination, id int) (*domain.GetAllWinnersCategoryResponse, error) {
	return w.repos.GetAll(page, id)
}

func (w *WinnerService) GetAllMembers(page domain.Pagination, id int) (*domain.GetAllWinnersCategoryResponse, error) {
	return w.repos.GetAllMembers(page, id)
}

func (w *WinnerService) GetAllDays(page domain.Pagination, month int) (*domain.GetAllDaysResponse, error) {
	return w.repos.GetAllDays(page, month)
}
func (w *WinnerService) GetAllMonths(page domain.Pagination) (*domain.GetAllDaysResponse, error) {
	return w.repos.GetAllMonths(page)
}
