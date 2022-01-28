package repository

import (
	"HundredToFive/internal/domain"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"reflect"
	"strings"
)

type ConditionRepos struct {
	db *sqlx.DB
}

func NewConditionRepos(db *sqlx.DB) *ConditionRepos {
	return &ConditionRepos{db: db}
}

func (c *ConditionRepos) Create(con domain.Condition) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s(caption,text,language_type) VALUES($1,$2,$3) RETURNING id", conditions)

	err := c.db.QueryRowx(query, con.Caption, con.Text, con.LanguageType).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("repository.Create: %w", err)
	}
	return id, nil
}

func (c *ConditionRepos) GetAll(page domain.Pagination, lang string) (*domain.GetAllConditionCategoryResponse, error) {
	var (
		setValues string
		count     int
	)
	if lang != "" {
		setValues = fmt.Sprintf("WHERE language_type = '%s'", lang)
	}

	queryCount := fmt.Sprintf(
		`SELECT COUNT(*) FROM %s %s`, conditions, setValues)

	err := c.db.QueryRowx(queryCount).Scan(&count)

	offset, pagesCount := calculatePagination(&page, count)

	inp := make([]*domain.Condition, 0, page.Limit)

	query := fmt.Sprintf("SELECT * FROM %s %s ORDER BY id ASC LIMIT $1 OFFSET $2", conditions, setValues)

	err = c.db.Select(&inp, query, page.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("repository.GetAll: %w", err)
	}

	pages := domain.PaginationPage{
		Page:  page.Page,
		Pages: pagesCount,
		Count: count,
	}

	con := domain.GetAllConditionCategoryResponse{
		Data:     inp,
		PageInfo: pages,
	}
	return &con, nil
}

func (c *ConditionRepos) GetById(id int) (domain.Condition, error) {
	var con domain.Condition
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", conditions)
	err := c.db.Get(&con, query, id)
	if err != nil {
		return domain.Condition{}, fmt.Errorf("repository.GetById: %w", domain.ErrNotFound)
	}
	return con, nil
}

func (c *ConditionRepos) Update(id int, inp domain.Condition) error {
	setValues := make([]string, 0, reflect.TypeOf(domain.Condition{}).NumField())

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

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=%d", conditions, setQuery, id)

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

func (c *ConditionRepos) Delete(id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", conditions)
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
