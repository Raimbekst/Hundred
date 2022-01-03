package repository

import (
	"HundredToFive/internal/domain"
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"
)

type CheckRepository struct {
	db *sqlx.DB
}

func NewCheckRepository(db *sqlx.DB) *CheckRepository {
	return &CheckRepository{db: db}
}

func (c *CheckRepository) Create(info domain.CheckInfo) (int, error) {
	tx := c.db.MustBegin()

	query := fmt.Sprintf(
		`INSERT INTO 
					%s
				(user_id,partner_id,check_amount,check_date)
					VALUES($1,$2,$3, to_timestamp($4) at time zone 'GMT')
				RETURNING id `, checks)

	var id int

	err := tx.QueryRowx(query, info.UserId, info.PartnerId, info.CheckAmount, info.CheckDate).Scan(&id)

	if err != nil {

		if txErr := tx.Rollback(); txErr != nil {
			return 0, fmt.Errorf("repository.Create: %w", txErr)
		}

		return 0, fmt.Errorf("repository.Create: %w", err)
	}
	fmt.Println(id)

	var imageList domain.CheckImage

	if info.CheckList != nil {

		for _, value := range info.CheckList {
			imageList.CheckId = id
			imageList.CheckImage = value

			info.CheckImage = append(info.CheckImage, imageList)
		}

		for _, value := range info.CheckImage {
			fmt.Println(value)
		}

		queryCheckImg := fmt.Sprintf(
			`INSERT INTO 
					%s
				(check_id,check_image) 
					VALUES (:check_id,:check_image)`,
			checkImages,
		)

		rows, err := tx.NamedExec(queryCheckImg, info.CheckImage)
		if err != nil {
			if txErr := tx.Rollback(); txErr != nil {
				return 0, fmt.Errorf("repository.Create: %w", txErr)
			}
			return 0, fmt.Errorf("repository.Create: %w", err)
		}

		affected, err := rows.RowsAffected()
		if err != nil {
			if txErr := tx.Rollback(); txErr != nil {
				return 0, fmt.Errorf("repository.Create: %w", txErr)
			}
			return 0, fmt.Errorf("repository.Create: %w", err)
		}

		if affected == 0 {
			if txErr := tx.Rollback(); txErr != nil {
				return 0, fmt.Errorf("repository.Create: %w", txErr)
			}
			return 0, fmt.Errorf("repository.Create: %w", sql.ErrNoRows)
		}

	}

	if txErr := tx.Commit(); txErr != nil {
		return 0, fmt.Errorf("repository.Create: %w", txErr)
	}

	return id, nil

}

func (c *CheckRepository) GetAll(ctx *fiber.Ctx, page domain.Pagination, filter domain.FilterForCheck) (*domain.GetAllChecksCategoryResponse, error) {
	var (
		err            error
		count          int
		whereClause    string
		forCheckValues []string
		setValues      string
	)

	url := ctx.BaseURL()

	if filter.EndRegisteredDate == 0 {
		filter.EndRegisteredDate = float64(time.Now().Unix())
	}
	if filter.EndCheckDate == 0 {
		filter.EndCheckDate = float64(time.Now().Unix())
	}

	if filter.MoneyAmount != 0 {
		forCheckValues = append(forCheckValues, fmt.Sprintf("check_amount <= %d", filter.MoneyAmount))
	}

	if filter.PartnerId != 0 {
		forCheckValues = append(forCheckValues, fmt.Sprintf("partner_id = %d", filter.PartnerId))
	}

	forCheckValues = append(forCheckValues, fmt.Sprintf("c.registered_at between to_timestamp(%f)::timestamp and to_timestamp(%f)::timestamp", filter.StartRegisteredDate, filter.EndRegisteredDate))

	forCheckValues = append(forCheckValues, fmt.Sprintf("check_date between to_timestamp(%f) and to_timestamp(%f) at time zone 'GMT'", filter.StartCheckDate, filter.EndCheckDate))

	forCheckValues = append(forCheckValues, fmt.Sprintf("is_winner = %v", false))

	whereClause = strings.Join(forCheckValues, " AND ")

	if whereClause != "" {
		setValues = "WHERE " + whereClause
	}

	queryCount := fmt.Sprintf(
		`SELECT COUNT(*) FROM %s c %s`, checks, setValues)

	err = c.db.QueryRowx(queryCount).Scan(&count)

	if err != nil {
		return nil, fmt.Errorf("repository.GetAll: %w", err)
	}

	offset, pagesCount := calculatePagination(&page, count)

	checkImage := make([]*domain.UserChecks, 0, page.Limit)

	query := fmt.Sprintf(

		`SELECT 
					c.user_id, 
					extract(epoch from c.check_date::timestamp at time zone 'GMT') "check_date",
					c.id,
					extract(epoch from c.registered_at) "registered_at", 
					c.check_amount,
					us.user_name,
					us.phone_number, 
					us.is_blocked, 
					partner.partner_name
				FROM 
					%s c 
				INNER JOIN 
					%s us
				ON 
					c.user_id = us.id
				INNER JOIN
					%s partner
				ON 
					c.partner_id = partner.id
				%s
				   ORDER BY 
					id 
				ASC 
					LIMIT $1 OFFSET $2`,
		checks, usersTable, partners, setValues)

	err = c.db.Select(&checkImage, query, page.Limit, offset)

	if err != nil {
		return nil, fmt.Errorf("repository.GetAll: %w", err)
	}

	for _, value := range checkImage {
		queryImg := fmt.Sprintf("SELECT * FROM %s WHERE check_id = $1", checkImages)
		err := c.db.Select(&value.CheckImage, queryImg, value.Id)

		if err != nil {
			return nil, fmt.Errorf("repository.GetAll: %w", err)
		}

		for _, val := range value.CheckImage {
			val.CheckImage = url + "/" + "media/" + val.CheckImage
		}
	}

	pages := domain.PaginationPage{
		Page:  page.Page,
		Pages: pagesCount,
		Count: count,
	}

	checks := domain.GetAllChecksCategoryResponse{
		Data:     checkImage,
		PageInfo: pages,
	}

	return &checks, nil
}

func (c *CheckRepository) GetById(ctx *fiber.Ctx, id int) (*domain.UserChecks, error) {
	url := ctx.BaseURL()
	var checkImage domain.UserChecks

	query := fmt.Sprintf(
		`SELECT 
					c.user_id, 
					extract(epoch from c.check_date::timestamp at time zone 'GMT') "check_date",
					c.id,
					extract(epoch from c.registered_at) "registered_at", 
					c.check_amount,
					us.user_name,
					us.phone_number, 
					us.is_blocked, 
					partner.partner_name
				FROM 
					%s c 
				INNER JOIN 
					%s us
				ON 
					c.user_id = us.id
				INNER JOIN
					%s partner
				ON 
					c.partner_id = partner.id WHERE c.id = $1`,
		checks, usersTable, partners)

	err := c.db.Get(&checkImage, query, id)

	if err != nil {
		return nil, fmt.Errorf("repository.GetAll: %w", domain.ErrNotFound)
	}

	queryImg := fmt.Sprintf("SELECT * FROM %s WHERE check_id = $1", checkImages)
	err = c.db.Select(&checkImage.CheckImage, queryImg, checkImage.Id)

	if err != nil {
		return nil, fmt.Errorf("repository.GetAll: %w", err)
	}
	for _, val := range checkImage.CheckImage {
		val.CheckImage = url + "/" + "media/" + val.CheckImage
	}

	return &checkImage, nil

}
