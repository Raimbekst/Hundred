package service

import (
	"HundredToFive/internal/domain"
	"HundredToFive/internal/repository"
)

type FaqService struct {
	repos repository.Faq
}

func NewFaqService(repos repository.Faq) *FaqService {
	return &FaqService{repos: repos}
}

func (f *FaqService) Create(faq domain.Faq) (int, error) {
	return f.repos.Create(faq)
}

func (f *FaqService) GetAll(page domain.Pagination, lang string) (*domain.GetAllFaqsCategoryResponse, error) {
	return f.repos.GetAll(page, lang)
}

func (f *FaqService) GetById(id int) (domain.Faq, error) {
	return f.repos.GetById(id)
}

func (f *FaqService) Update(id int, inp domain.Faq) error {
	return f.repos.Update(id, inp)
}

func (f *FaqService) Delete(id int) error {
	return f.repos.Delete(id)
}

func (f *FaqService) CreateDesc(desc domain.Description) (int, error) {
	return f.repos.CreateDesc(desc)
}

func (f *FaqService) GetAllDesc(page domain.Pagination, lang string) (*domain.GetAllDescCategoryResponse, error) {
	return f.repos.GetAllDesc(page, lang)
}

func (f *FaqService) GetDescById(id int) (domain.Description, error) {
	return f.repos.GetDescById(id)
}

func (f *FaqService) UpdateDesc(id int, inp domain.Description) error {
	return f.repos.UpdateDesc(id, inp)
}

func (f *FaqService) DeleteDesc(id int) error {
	return f.repos.DeleteDesc(id)
}
