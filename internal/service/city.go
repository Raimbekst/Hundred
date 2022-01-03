package service

import (
	"HundredToFive/internal/domain"
	"HundredToFive/internal/repository"
)

type CityService struct {
	repos repository.City
}

func NewCityService(repos repository.City) *CityService {
	return &CityService{repos: repos}
}

func (c *CityService) Create(city domain.City) (int, error) {
	return c.repos.Create(city)
}
func (c *CityService) GetAll(page domain.Pagination) (domain.GetAllCityCategoryResponse, error) {
	return c.repos.GetAll(page)
}

func (c *CityService) GetById(id int) (domain.City, error) {
	return c.repos.GetById(id)
}
func (c *CityService) Update(id int, inp domain.City) error {
	return c.repos.Update(id, inp)
}
func (c *CityService) Delete(id int) error {
	return c.repos.Delete(id)
}
