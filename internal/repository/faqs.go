package repository

import (
	"HundredToFive/internal/domain"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"reflect"
	"strings"
)

type FaqsRepos struct {
	db *sqlx.DB
}

func (f *FaqsRepos) CreateDesc(desc domain.Description) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s(caption,text,language_type) VALUES($1,$2,$3) RETURNING id", descriptions)

	err := f.db.QueryRowx(query, desc.Caption, desc.Text, desc.LanguageType).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("repository.Create: %w", err)
	}
	return id, nil
}

func (f *FaqsRepos) GetAllDesc(page domain.Pagination, lang string) (*domain.GetAllDescCategoryResponse, error) {
	var (
		setValues string
		count     int
	)
	if lang != "" {
		setValues = fmt.Sprintf("WHERE language_type = '%s'", lang)
	}

	queryCount := fmt.Sprintf(
		`SELECT COUNT(*) FROM %s %s`, descriptions, setValues)

	err := f.db.QueryRowx(queryCount).Scan(&count)

	offset, pagesCount := calculatePagination(&page, count)

	inp := make([]*domain.Description, 0, page.Limit)

	query := fmt.Sprintf("SELECT * FROM %s %s ORDER BY id ASC LIMIT $1 OFFSET $2", descriptions, setValues)

	err = f.db.Select(&inp, query, page.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("repository.GetAll: %w", err)
	}

	pages := domain.PaginationPage{
		Page:  page.Page,
		Pages: pagesCount,
		Count: count,
	}

	desc := domain.GetAllDescCategoryResponse{
		Data:     inp,
		PageInfo: pages,
	}
	return &desc, nil
}

func (f *FaqsRepos) GetDescById(id int) (domain.Description, error) {
	var faq domain.Description
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", descriptions)
	err := f.db.Get(&faq, query, id)
	if err != nil {
		return domain.Description{}, fmt.Errorf("repository.GetById: %w", domain.ErrNotFound)
	}
	return faq, nil
}

func (f *FaqsRepos) UpdateDesc(id int, inp domain.Description) error {
	setValues := make([]string, 0, reflect.TypeOf(domain.Description{}).NumField())

	if inp.Caption != "" {
		setValues = append(setValues, fmt.Sprintf("caption=:caption"))
	}
	if inp.Text != "" {
		setValues = append(setValues, fmt.Sprintf("text=:text"))
	}

	setQuery := strings.Join(setValues, ", ")

	if setQuery == "" {
		return errors.New("empty body")
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=%d", descriptions, setQuery, id)

	result, err := f.db.NamedExec(query, inp)

	if err != nil {
		return fmt.Errorf("repository.Update: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("repository.Update: %w", err)
	}

	if affected == 0 {
		return fmt.Errorf("repository.Update: %w", domain.ErrNotFound)
	}

	return nil
}

func (f *FaqsRepos) DeleteDesc(id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", descriptions)
	result, err := f.db.Exec(query, id)
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("repository.Delete: %w", err)
	}

	if affected == 0 {
		return fmt.Errorf("repository.Delete: %w", domain.ErrNotFound)
	}

	return nil
}

func NewFaqsRepos(db *sqlx.DB) *FaqsRepos {
	return &FaqsRepos{db: db}
}

func (f *FaqsRepos) Create(faq domain.Faq) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s(question,answer,language_type) VALUES($1,$2,$3) RETURNING id", faqs)

	err := f.db.QueryRowx(query, faq.Question, faq.Answer, faq.LanguageType).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("repository.Create: %w", err)
	}
	return id, nil
}

func (f *FaqsRepos) GetAll(page domain.Pagination, lang string) (*domain.GetAllFaqsCategoryResponse, error) {
	var (
		setValues string
		count     int
	)
	if lang != "" {
		setValues = fmt.Sprintf("WHERE language_type = '%s'", lang)
	}

	queryCount := fmt.Sprintf(
		`SELECT COUNT(*) FROM %s %s`, faqs, setValues)

	err := f.db.QueryRowx(queryCount).Scan(&count)

	offset, pagesCount := calculatePagination(&page, count)

	inp := make([]*domain.Faq, 0, page.Limit)

	query := fmt.Sprintf("SELECT * FROM %s %s ORDER BY id ASC LIMIT $1 OFFSET $2", faqs, setValues)
	err = f.db.Select(&inp, query, page.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("repository.GetAll: %w", err)
	}

	pages := domain.PaginationPage{
		Page:  page.Page,
		Pages: pagesCount,
		Count: count,
	}

	faq := domain.GetAllFaqsCategoryResponse{
		Data:     inp,
		PageInfo: pages,
	}
	return &faq, nil
}

func (f FaqsRepos) GetById(id int) (domain.Faq, error) {
	var faq domain.Faq
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", faqs)
	err := f.db.Get(&faq, query, id)
	if err != nil {
		return domain.Faq{}, fmt.Errorf("repository.GetById: %w", domain.ErrNotFound)
	}
	return faq, nil
}

func (f *FaqsRepos) Update(id int, inp domain.Faq) error {
	setValues := make([]string, 0, reflect.TypeOf(domain.Faq{}).NumField())

	if inp.Answer != "" {
		setValues = append(setValues, fmt.Sprintf("answer=:answer"))
	}
	if inp.Question != "" {
		setValues = append(setValues, fmt.Sprintf("question=:question"))
	}

	setQuery := strings.Join(setValues, ", ")

	if setQuery == "" {
		return errors.New("empty body")
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=%d", faqs, setQuery, id)

	result, err := f.db.NamedExec(query, inp)

	if err != nil {
		return fmt.Errorf("repository.Update: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("repository.Update: %w", err)
	}

	if affected == 0 {
		return fmt.Errorf("repository.Update: %w", domain.ErrNotFound)
	}

	return nil
}

func (f *FaqsRepos) Delete(id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", faqs)
	result, err := f.db.Exec(query, id)
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("repository.Delete: %w", err)
	}

	if affected == 0 {
		return fmt.Errorf("repository.Delete: %w", domain.ErrNotFound)
	}

	return nil
}
