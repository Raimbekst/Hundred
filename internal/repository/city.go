package repository

import (
	"HundredToFive/internal/domain"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"reflect"
	"strings"
)

type CityRepos struct {
	db *sqlx.DB
}

func NewCityRepos(db *sqlx.DB) *CityRepos {
	return &CityRepos{db: db}
}

func (c *CityRepos) Create(city domain.City) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s(name) VALUES($1) RETURNING id", cities)
	err := c.db.QueryRowx(query, city.Name).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("repository.Create: %w", err)
	}
	return id, nil
}
func (c *CityRepos) GetAll(page domain.Pagination) (domain.GetAllCityCategoryResponse, error) {
	count, err := countPage(c.db, cities)
	if err != nil {
		return domain.GetAllCityCategoryResponse{}, fmt.Errorf("repository.GetAll: %w", err)
	}

	offset, pagesCount := calculatePagination(&page, count)

	inp := make([]*domain.City, 0, page.Limit)

	query := fmt.Sprintf("SELECT * FROM %s ORDER BY id ASC LIMIT $1 OFFSET $2", cities)
	err = c.db.Select(&inp, query, page.Limit, offset)
	if err != nil {
		return domain.GetAllCityCategoryResponse{}, fmt.Errorf("repository.GetAll: %w", err)
	}

	pages := domain.PaginationPage{
		Page:  page.Page,
		Pages: pagesCount,
		Count: count,
	}

	city := domain.GetAllCityCategoryResponse{
		Data:     inp,
		PageInfo: pages,
	}
	return city, nil
}

func (c *CityRepos) GetById(id int) (domain.City, error) {
	var city domain.City
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", cities)
	err := c.db.Get(&city, query, id)
	if err != nil {
		return domain.City{}, fmt.Errorf("repository.GetById: %w", domain.ErrNotFound)
	}
	return city, nil
}
func (c *CityRepos) Update(id int, inp domain.City) error {
	setValues := make([]string, 0, reflect.TypeOf(domain.City{}).NumField())

	if inp.Name != "" {
		setValues = append(setValues, fmt.Sprintf("name=:name"))
	}

	setQuery := strings.Join(setValues, ", ")

	if setQuery == "" {
		return errors.New("empty body")
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=%d", cities, setQuery, id)

	result, err := c.db.NamedExec(query, inp)

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
func (c *CityRepos) Delete(id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", cities)
	result, err := c.db.Exec(query, id)
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("repository.Delete: %w", err)
	}

	if affected == 0 {
		return fmt.Errorf("repository.Delete: %w", domain.ErrNotFound)
	}

	return nil
}
