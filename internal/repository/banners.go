package repository

import (
	"HundredToFive/internal/domain"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"reflect"
	"strings"
)

type BannerRepos struct {
	db *sqlx.DB
}

func NewBannerRepos(db *sqlx.DB) *BannerRepos {
	return &BannerRepos{db: db}
}

func (b *BannerRepos) Create(banner domain.Banner) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s(name,status,image,iframe,language_type) VALUES($1,$2,$3,$4,$5) RETURNING id", banners)
	err := b.db.QueryRowx(query, banner.Name, banner.Status, banner.Image, banner.Iframe, banner.LanguageType).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("repository.Create: %w", err)
	}
	return id, nil
}

func (b *BannerRepos) GetAll(page domain.Pagination, status int, lang string) (*domain.GetAllBannersCategoryResponse, error) {
	var (
		setValues      string
		count          int
		forCheckValues []string
		whereClause    string
	)
	if status != 0 {
		forCheckValues = append(forCheckValues, fmt.Sprintf("status = %d", status))
	}
	if lang != "" {
		forCheckValues = append(forCheckValues, fmt.Sprintf("language_type = '%s'", lang))

	}
	whereClause = strings.Join(forCheckValues, " AND ")

	if whereClause != "" {
		setValues = "WHERE " + whereClause
	}

	queryCount := fmt.Sprintf(
		`SELECT COUNT(*) FROM %s  %s`, banners, setValues)

	err := b.db.QueryRowx(queryCount).Scan(&count)

	if err != nil {
		return nil, fmt.Errorf("repository.GetAll: %w", err)
	}
	offset, pagesCount := calculatePagination(&page, count)

	inp := make([]*domain.Banner, 0, page.Limit)

	query := fmt.Sprintf("SELECT * FROM %s %s ORDER BY id ASC LIMIT $1 OFFSET $2", banners, setValues)
	err = b.db.Select(&inp, query, page.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("repository.GetAll: %w", err)
	}

	pages := domain.PaginationPage{
		Page:  page.Page,
		Pages: pagesCount,
		Count: count,
	}

	banner := domain.GetAllBannersCategoryResponse{
		Data:     inp,
		PageInfo: pages,
	}
	return &banner, nil
}

func (b *BannerRepos) GetById(id int) (domain.Banner, error) {
	var banner domain.Banner
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", banners)
	err := b.db.Get(&banner, query, id)
	if err != nil {
		return domain.Banner{}, fmt.Errorf("repository.GetById: %w", domain.ErrNotFound)
	}
	return banner, nil
}

func (b *BannerRepos) Update(id int, inp domain.Banner) (string, error) {

	var imageInput domain.Banner

	querySelect := fmt.Sprintf("SELECT image FROM %s WHERE id = $1", banners)

	err := b.db.Get(&imageInput, querySelect, id)

	if err != nil {
		return "", fmt.Errorf("repository.Update: %w", err)
	}

	setValues := make([]string, 0, reflect.TypeOf(domain.Banner{}).NumField())

	if inp.Name != "" {
		setValues = append(setValues, fmt.Sprintf("name=:name"))
	}

	if inp.Status != 0 {
		setValues = append(setValues, fmt.Sprintf("status=:status"))
	}

	if inp.Image != "" {
		setValues = append(setValues, fmt.Sprintf("image=:image"))
	}

	if inp.Iframe != "" {
		setValues = append(setValues, fmt.Sprintf("iframe=:iframe"))
	}

	setQuery := strings.Join(setValues, ", ")

	if setQuery == "" {
		return "", errors.New("empty body")
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=%d", banners, setQuery, id)

	rows, err := b.db.NamedExec(query, inp)

	if err != nil {
		return "", fmt.Errorf("repository.Update: %w", err)
	}
	affected, err := rows.RowsAffected()
	if err != nil {
		return "", fmt.Errorf("repository.Update: %w", err)
	}
	if affected == 0 {
		return "", fmt.Errorf("repository.Update: %w", err)
	}
	return imageInput.Image, nil
}

func (b *BannerRepos) Delete(id int) (string, error) {
	var image string
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1 RETURNING image", banners)
	err := b.db.QueryRow(query, id).Scan(&image)
	if err != nil {
		return "", fmt.Errorf("repository.Delete: %w", domain.ErrNotFound)
	}
	return image, nil
}
