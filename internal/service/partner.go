package service

import (
	"HundredToFive/internal/domain"
	"HundredToFive/internal/repository"
	"HundredToFive/pkg/media"
	"fmt"
	"github.com/xuri/excelize/v2"
	"strconv"
)

type PartnerService struct {
	repos repository.Partner
}

func NewPartnerService(repos repository.Partner) *PartnerService {
	return &PartnerService{repos: repos}
}

func (p *PartnerService) Create(partner domain.Partner) (int, error) {
	return p.repos.Create(partner)
}
func (p *PartnerService) GetAll(page domain.Pagination, status int) (*domain.GetAllPartnersCategoryResponse, error) {
	return p.repos.GetAll(page, status)
}

func (p *PartnerService) DownloadPartners(file *excelize.File, url string) (*excelize.File, error) {
	page := domain.Pagination{}
	status := 0
	list, err := p.repos.GetAll(page, status)

	if err != nil {
		return nil, fmt.Errorf("service.DownloadUsers: %w", err)
	}
	id := 2

	for _, value := range list.Data {
		file.SetCellValue("Sheet1", "A"+strconv.Itoa(id), value.Id)
		file.SetCellValue("Sheet1", "B"+strconv.Itoa(id), value.PartnerName)
		file.SetCellValue("Sheet1", "C"+strconv.Itoa(id), value.LinkWebsite)
		file.SetCellValue("Sheet1", "D"+strconv.Itoa(id), value.Status)
		file.SetCellValue("Sheet1", "E"+strconv.Itoa(id), value.Reference)
		file.SetCellValue("Sheet1", "F"+strconv.Itoa(id), url+"/"+"media/"+value.Logo)
		file.SetCellValue("Sheet1", "G"+strconv.Itoa(id), url+"/"+"media/"+value.Banner)
		file.SetCellValue("Sheet1", "H"+strconv.Itoa(id), url+"/"+"media/"+value.BannerKz)
		file.SetCellValue("Sheet1", "I"+strconv.Itoa(id), value.StartPartnership)
		file.SetCellValue("Sheet1", "J"+strconv.Itoa(id), value.EndPartnership)
		id = id + 1
	}
	return file, nil
}

func (p *PartnerService) GetById(id int) (domain.Partner, error) {
	return p.repos.GetById(id)
}
func (p *PartnerService) Update(id int, inp domain.Partner) error {
	img, err := p.repos.Update(id, inp)
	for i, _ := range img {
		if img[i] != "" {
			err = media.DeleteImage(img[i])
			if err != nil {
				return fmt.Errorf("service.Delete: %w", err)
			}
		}
	}
	return nil
}
func (p *PartnerService) Delete(id int) error {
	img, err := p.repos.Delete(id)
	if err != nil {
		return fmt.Errorf("service.Delete: %w", err)

	}
	for i, _ := range img {
		if img[i] != "" {
			err = media.DeleteImage(img[i])
			if err != nil {
				return fmt.Errorf("service.Delete: %w", err)
			}
		}
	}
	return nil
}
