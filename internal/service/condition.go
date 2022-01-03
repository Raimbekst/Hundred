package service

import (
	"HundredToFive/internal/domain"
	"HundredToFive/internal/repository"
)

type ConditionService struct {
	repos repository.Condition
}

func NewConditionService(repos repository.Condition) *ConditionService {
	return &ConditionService{repos: repos}
}

func (c *ConditionService) Create(con domain.Condition) (int, error) {
	return c.repos.Create(con)
}

func (c *ConditionService) GetAll(page domain.Pagination) (*domain.GetAllConditionCategoryResponse, error) {
	return c.repos.GetAll(page)
}

func (c *ConditionService) GetById(id int) (domain.Condition, error) {
	return c.repos.GetById(id)
}

func (c *ConditionService) Update(id int, inp domain.Condition) error {
	return c.repos.Update(id, inp)
}

func (c *ConditionService) Delete(id int) error {
	return c.repos.Delete(id)
}
