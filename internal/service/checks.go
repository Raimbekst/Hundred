package service

import (
	"HundredToFive/internal/domain"
	"HundredToFive/internal/repository"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
	"strconv"
	"time"
)

type CheckService struct {
	repos repository.Check
}

func NewCheckService(repos repository.Check) *CheckService {
	return &CheckService{repos: repos}
}

func (c *CheckService) Create(info domain.CheckInfo) (int, error) {
	id, err := c.repos.Create(info)
	if err != nil {
		return 0, fmt.Errorf("service.Create: %w", err)
	}
	return id, nil
}

func (c *CheckService) GetAll(ctx *fiber.Ctx, page domain.Pagination, filter domain.FilterForCheck) (*domain.GetAllChecksCategoryResponse, error) {
	list, err := c.repos.GetAll(ctx, page, filter)
	if err != nil {
		return nil, fmt.Errorf("service.GetAll: %w", err)
	}
	return list, nil
}

func (c *CheckService) GetById(ctx *fiber.Ctx, id int) (*domain.UserChecks, error) {
	return c.repos.GetById(ctx, id)
}

func (c *CheckService) DownloadChecks(ctx *fiber.Ctx, file *excelize.File, filter domain.FilterForCheck) (*excelize.File, error) {
	page := domain.Pagination{}

	list, err := c.repos.GetAll(ctx, page, filter)

	if err != nil {
		return nil, fmt.Errorf("service.DownloadChecks: %w", err)
	}
	id := 2

	for _, value := range list.Data {
		unixTime := time.Unix(int64(value.RegisteredAt), 0)
		dateCheck := time.Unix(int64(value.CheckDate), 0)
		file.SetCellValue("Sheet1", "A"+strconv.Itoa(id), value.Id)
		file.SetCellValue("Sheet1", "B"+strconv.Itoa(id), value.UserId)
		file.SetCellValue("Sheet1", "C"+strconv.Itoa(id), value.UserName)
		file.SetCellValue("Sheet1", "D"+strconv.Itoa(id), value.PhoneNumber)
		file.SetCellValue("Sheet1", "E"+strconv.Itoa(id), value.PartnerName)
		file.SetCellValue("Sheet1", "F"+strconv.Itoa(id), unixTime)
		file.SetCellValue("Sheet1", "G"+strconv.Itoa(id), dateCheck)
		file.SetCellValue("Sheet1", "H"+strconv.Itoa(id), value.CheckAmount)
		file.SetCellValue("Sheet1", "I"+strconv.Itoa(id), value.CheckImage)
		id = id + 1
	}
	return file, nil
}
