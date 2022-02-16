package repository

import (
	"HundredToFive/internal/domain"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"reflect"
	"strings"
)

type PartnerRepos struct {
	db *sqlx.DB
}

func NewPartnerRepos(db *sqlx.DB) *PartnerRepos {
	return &PartnerRepos{db: db}
}

func (p *PartnerRepos) Create(partner domain.Partner) (int, error) {
	var id int
	query := fmt.Sprintf(
		`INSERT INTO 
				%s
						(partner_name,position,logo,link_website,banner,banner_kz,status,start_partnership,end_partnership,partner_package,reference) 
				VALUES
						($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) 
				RETURNING id`, partners)

	err := p.db.QueryRowx(query, partner.PartnerName, partner.Position, partner.Logo, partner.LinkWebsite, partner.Banner, partner.BannerKz, partner.Status, partner.StartPartnership, partner.EndPartnership, partner.PartnerPackage, partner.Reference).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("repository.Create: %w", err)
	}
	return id, nil
}

func (p *PartnerRepos) GetAll(page domain.Pagination, status int) (*domain.GetAllPartnersCategoryResponse, error) {
	var (
		setValues string
		count     int
	)
	if status != 0 {
		setValues = fmt.Sprintf("WHERE status = %d", status)
	}

	queryCount := fmt.Sprintf(
		`SELECT COUNT(*) FROM %s  %s`, partners, setValues)

	err := p.db.QueryRowx(queryCount).Scan(&count)

	if err != nil {
		return nil, fmt.Errorf("repository.GetAll: %w", err)
	}

	offset, pagesCount := calculatePagination(&page, count)

	inp := make([]*domain.Partner, 0, page.Limit)

	query := fmt.Sprintf("SELECT * FROM %s %s ORDER BY position ASC LIMIT $1 OFFSET $2", partners, setValues)

	err = p.db.Select(&inp, query, page.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("repository.GetAll: %w", err)
	}

	pages := domain.PaginationPage{
		Page:  page.Page,
		Pages: pagesCount,
		Count: count,
	}

	partner := domain.GetAllPartnersCategoryResponse{
		Data:     inp,
		PageInfo: pages,
	}
	return &partner, nil
}

func (p PartnerRepos) GetById(id int) (domain.Partner, error) {
	var partner domain.Partner
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", partners)
	err := p.db.Get(&partner, query, id)
	if err != nil {
		return domain.Partner{}, fmt.Errorf("repository.GetById: %w", domain.ErrNotFound)
	}
	return partner, nil
}

func (p *PartnerRepos) Update(id int, inp domain.UpdatePartner) ([]string, error) {

	var images []string

	part := domain.Partner{
		Position:         inp.Position,
		PartnerName:      inp.PartnerName,
		LinkWebsite:      inp.LinkWebsite,
		Status:           inp.Status,
		StartPartnership: inp.StartPartnership,
		EndPartnership:   inp.EndPartnership,
		PartnerPackage:   inp.PartnerPackage,
		Reference:        inp.Reference,
	}

	var imageInput domain.Partner
	setValues := make([]string, 0, reflect.TypeOf(domain.UpdatePartner{}).NumField())

	if inp.PartnerName != "" {
		setValues = append(setValues, fmt.Sprintf("partner_name=:partner_name"))

	}
	if inp.Logo != nil {
		setValues = append(setValues, fmt.Sprintf("logo=:logo"))
		images = append(images, "logo")
		part.Logo = *inp.Logo
	}
	if inp.LinkWebsite != "" {
		setValues = append(setValues, fmt.Sprintf("link_website=:link_website"))
	}

	if inp.Position != 0 {
		setValues = append(setValues, fmt.Sprintf("position=:position"))
	}

	if inp.Banner != nil {
		setValues = append(setValues, fmt.Sprintf("banner=:banner"))
		images = append(images, "banner")
		part.Banner = *inp.Banner
	}
	if inp.BannerKz != nil {
		setValues = append(setValues, fmt.Sprintf("banner_kz=:banner_kz"))
		images = append(images, "banner_kz")
		part.BannerKz = *inp.BannerKz
	}
	if inp.Status != 0 {
		setValues = append(setValues, fmt.Sprintf("status=:status"))
	}
	if inp.StartPartnership != "" {
		setValues = append(setValues, fmt.Sprintf("start_partnership=:start_partnership"))
	}

	if inp.EndPartnership != "" {
		setValues = append(setValues, fmt.Sprintf("end_partnership=:end_partnership"))
	}
	if inp.PartnerPackage != "" {
		setValues = append(setValues, fmt.Sprintf("partner_package=:partner_package"))
	}
	if inp.Reference != "" {
		setValues = append(setValues, fmt.Sprintf("reference=:reference"))
	}

	imageString := strings.Join(images, ", ")

	querySelect := fmt.Sprintf("SELECT %s FROM %s WHERE id = $1", imageString, partners)

	err := p.db.Get(&imageInput, querySelect, id)

	if err != nil {
		return nil, fmt.Errorf("repository.Update: %w", err)
	}

	images = nil
	images = append(images, imageInput.Logo, imageInput.Banner, imageInput.BannerKz)

	setQuery := strings.Join(setValues, ", ")

	if setQuery == "" {
		return nil, errors.New("empty body")
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=%d", partners, setQuery, id)

	result, err := p.db.NamedExec(query, part)

	if err != nil {
		return nil, fmt.Errorf("repository.Update: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("repository.Update: %w", err)
	}

	if affected == 0 {
		return nil, fmt.Errorf("repository.Update: %w", domain.ErrNotFound)
	}
	return images, nil

}
func (p *PartnerRepos) Delete(id int) ([]string, error) {
	var (
		images   []string
		logo     string
		banner   string
		bannerKz string
	)
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1 RETURNING logo,banner,banner_kz", partners)
	err := p.db.QueryRow(query, id).Scan(&logo, &banner, &bannerKz)

	if err != nil {
		return nil, fmt.Errorf("repository.Delete: %w", domain.ErrNotFound)
	}

	images = append(images, logo, banner, bannerKz)
	return images, nil
}
