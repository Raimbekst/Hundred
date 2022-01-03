package service

import (
	"HundredToFive/internal/domain"
	"HundredToFive/internal/repository"
	"time"
)

type RaffleService struct {
	repos repository.Raffle
}

func NewRaffleService(repos repository.Raffle) *RaffleService {
	return &RaffleService{repos: repos}
}

func (r *RaffleService) Create(raffle domain.Raffle) (int, error) {

	status := map[int]string{1: planned, 2: finished, 3: unfinished}

	totalDate := raffle.RaffleDate + float64(raffle.RaffleTime)

	statusRaffle := status[1]

	if totalDate < float64(time.Now().Unix()) {
		statusRaffle = status[3]
	}

	raf := domain.Raffle{
		RaffleDate:    raffle.RaffleDate,
		RaffleTime:    raffle.RaffleTime,
		CheckCategory: raffle.CheckCategory,
		RaffleType:    raffle.RaffleType,
		Reference:     raffle.Reference,
		Status:        statusRaffle,
	}

	return r.repos.Create(raf)

}

func (r *RaffleService) GetAll(page domain.Pagination, filter domain.FilterForRaffles) (*domain.GetAllRaffleCategoryResponse, error) {
	return r.repos.GetAll(page, filter)
}

func (r *RaffleService) GetById(id int) (domain.Raffle, error) {
	return r.repos.GetById(id)
}
func (r *RaffleService) Update(id int, inp domain.Raffle) error {

	status := map[int]string{1: planned, 2: finished, 3: unfinished}

	totalDate := inp.RaffleDate + float64(inp.RaffleTime)

	statusRaffle := status[1]

	if totalDate < float64(time.Now().Unix()) {
		statusRaffle = status[3]
	}

	raf := domain.Raffle{
		RaffleDate:    inp.RaffleDate,
		RaffleTime:    inp.RaffleTime,
		CheckCategory: inp.CheckCategory,
		RaffleType:    inp.RaffleType,
		Reference:     inp.Reference,
		Status:        statusRaffle,
	}
	return r.repos.Update(id, raf)
}
func (r *RaffleService) Delete(id int) error {
	return r.repos.Delete(id)
}
