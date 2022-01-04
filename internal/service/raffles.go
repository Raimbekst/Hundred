package service

import (
	"HundredToFive/internal/domain"
	"HundredToFive/internal/repository"
	"fmt"
	"github.com/xuri/excelize/v2"
	"strconv"
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

func (r *RaffleService) DownloadRaffles(file *excelize.File) (*excelize.File, error) {
	page := domain.Pagination{}
	var filter domain.FilterForRaffles
	list, err := r.repos.GetAll(page, filter)

	if err != nil {
		return nil, fmt.Errorf("service.DownloadRaffles: %w", err)
	}
	id := 2
	wordNo := "нет"
	typeGame := map[int]string{1: "Ежедневный розыгрыш", 2: "Еженедельный розыгрыш", 3: "Ежемесячный розыгрыш"}
	for _, value := range list.Data {
		unixTime := time.Unix(int64(value.RaffleDate), 0)

		if value.UserName == nil {
			value.UserName = &wordNo
		}
		if value.PhoneNumber == nil {
			value.PhoneNumber = &wordNo
		}

		file.SetCellValue("Sheet1", "A"+strconv.Itoa(id), value.Id)
		file.SetCellValue("Sheet1", "B"+strconv.Itoa(id), *value.UserName)
		file.SetCellValue("Sheet1", "C"+strconv.Itoa(id), *value.PhoneNumber)
		file.SetCellValue("Sheet1", "D"+strconv.Itoa(id), unixTime)
		file.SetCellValue("Sheet1", "E"+strconv.Itoa(id), value.CheckCategory)
		file.SetCellValue("Sheet1", "F"+strconv.Itoa(id), typeGame[value.RaffleType])
		file.SetCellValue("Sheet1", "G"+strconv.Itoa(id), value.Status)
		id = id + 1
	}
	return file, nil
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
