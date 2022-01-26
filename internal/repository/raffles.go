package repository

import (
	"HundredToFive/internal/domain"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
)

type RaffleRepos struct {
	db *sqlx.DB
}

func NewRaffleRepos(db *sqlx.DB) *RaffleRepos {
	return &RaffleRepos{db: db}
}

func (c *RaffleRepos) Create(raffle domain.Raffle) (int, error) {
	var id int
	query := fmt.Sprintf(
		`INSERT INTO %s
						(raffle_date,raffle_time,check_category,raffle_type,status,reference)
					VALUES
							 (to_timestamp($1) at time zone 'GMT' ,$2, $3, $4, $5, $6) 
					RETURNING id`, raffles)

	err := c.db.QueryRowx(query, raffle.RaffleDate, raffle.RaffleTime, raffle.CheckCategory, raffle.RaffleType, raffle.Status, raffle.Reference).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("repository.Create: %w", err)
	}

	return id, nil
}

func (c *RaffleRepos) GetAll(page domain.Pagination, filter domain.FilterForRaffles) (*domain.GetAllRaffleCategoryResponse, error) {

	var (
		err             error
		count           int
		whereClause     string
		forRaffleValues []string
		setValues       string
	)
	isF := map[int]string{1: " = ", 2: " != "}

	if filter.RaffleType != 0 {
		forRaffleValues = append(forRaffleValues, fmt.Sprintf("r.raffle_type = %d", filter.RaffleType))
	}

	if filter.RaffleDate != 0 {
		forRaffleValues = append(forRaffleValues, fmt.Sprintf("r.raffle_date = to_timestamp(%f) at time zone 'GMT'", filter.RaffleDate))
	}

	if filter.RaffleTime != 0 {
		forRaffleValues = append(forRaffleValues, fmt.Sprintf("r.raffle_time = %d", filter.RaffleTime))
	}

	if filter.IsFinished != 0 {
		forRaffleValues = append(forRaffleValues, fmt.Sprintf("r.status %s '%s'", isF[filter.IsFinished], finished))
	}

	whereClause = strings.Join(forRaffleValues, " AND ")

	if whereClause != "" {
		setValues = "WHERE " + whereClause
	}
	fmt.Println(filter.RaffleDate)

	queryCount := fmt.Sprintf(
		`SELECT COUNT(*) FROM %s r %s`, raffles, setValues)

	err = c.db.QueryRowx(queryCount).Scan(&count)

	if err != nil {
		return nil, fmt.Errorf("repository.GetAll: %w", err)
	}

	offset, pagesCount := calculatePagination(&page, count)

	inp := make([]*domain.Raffle, 0, page.Limit)
	query := fmt.Sprintf(
		`select 
					r.id,
					r.raffle_time,
					r.raffle_type,
					r.status,
					r.check_category,
					r.reference,
					r.check_id,
					extract(epoch from raffle_date::timestamp at time zone 'GMT') "raffle_date",
					u.phone_number,
					u.user_name,
					u.id "user_id"
				from 
					%s r 
				left outer join 
				    %s ch 
				on 
					ch.id = r.check_id 
				left outer join 
					%s u 
				on 
					u.id = ch.user_id
				%s
				 ORDER BY r.id ASC LIMIT $1 OFFSET $2;`, raffles, checks, usersTable, setValues)

	err = c.db.Select(&inp, query, page.Limit, offset)

	if err != nil {
		return nil, fmt.Errorf("repository.GetAll: %w", err)
	}

	pages := domain.PaginationPage{
		Page:  page.Page,
		Pages: pagesCount,
		Count: count,
	}

	raffle := domain.GetAllRaffleCategoryResponse{
		Data:     inp,
		PageInfo: pages,
	}
	return &raffle, nil
}

func (c *RaffleRepos) GetById(id int) (domain.Raffle, error) {
	var raffle domain.Raffle
	query := fmt.Sprintf(
		`select 
					r.id,
					r.raffle_time,
					r.raffle_type,
					r.status,
					r.check_category,
					r.check_id,
					r.reference,
					extract(epoch from raffle_date::timestamp at time zone 'GMT') "raffle_date",
					u.phone_number,
					u.user_name,
					u.id "user_id"
				from 
					%s r 
				left outer join 
				    %s ch 
				on 
					ch.id = r.check_id 
				left outer join 
					%s u 
				on 
					u.id = ch.user_id
				WHERE r.id = $1`, raffles, checks, usersTable)
	err := c.db.Get(&raffle, query, id)
	if err != nil {
		return domain.Raffle{}, fmt.Errorf("repository.GetById: %w", domain.ErrNotFound)
	}
	return raffle, nil
}
func (c *RaffleRepos) Update(id int, inp domain.Raffle) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if inp.RaffleDate != 0 {
		setValues = append(setValues, fmt.Sprintf("raffle_date = to_timestamp($%d) at time zone 'GMT'", argId))
		args = append(args, inp.RaffleDate)
		argId++
	}

	if inp.CheckCategory != 0 {
		setValues = append(setValues, fmt.Sprintf("check_category=$%d", argId))
		args = append(args, inp.CheckCategory)
		argId++
	}

	if inp.RaffleType != 0 {
		setValues = append(setValues, fmt.Sprintf("raffle_type=$%d", argId))
		args = append(args, inp.RaffleType)
		argId++
	}

	if inp.Status != "" {
		setValues = append(setValues, fmt.Sprintf("status=$%d", argId))
		args = append(args, inp.Status)
		argId++
	}

	if inp.RaffleTime != 0 {
		setValues = append(setValues, fmt.Sprintf("raffle_time=$%d", argId))
		args = append(args, inp.RaffleTime)
		argId++
	}

	if inp.Reference != "" {
		setValues = append(setValues, fmt.Sprintf("reference=$%d", argId))
		args = append(args, inp.Reference)
		argId++
	}
	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf(`UPDATE %s SET %s WHERE id = $%d `, raffles, setQuery, argId)
	args = append(args, id)

	_, err := c.db.Exec(query, args...)
	return err

}

func (c *RaffleRepos) UpdateStatus(timeNow int64) error {
	query := fmt.Sprintf("UPDATE %s SET status = $1 WHERE raffle_date < to_timestamp($2) at time zone 'GMT' AND status = $3", raffles)

	_, err := c.db.Exec(query, notFinished, timeNow, planned)

	if err != nil {
		return fmt.Errorf("repository.UpdateStatus:%w", err)
	}
	return nil
}

func (c *RaffleRepos) Delete(id int) error {
	var checkId *int
	tx := c.db.MustBegin()

	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1 RETURNING check_id", raffles)
	err := tx.QueryRowx(query, id).Scan(&checkId)

	if err != nil {
		if txErr := tx.Rollback(); txErr != nil {
			return fmt.Errorf("repository.Delete: %w", txErr)
		}
		return fmt.Errorf("repository.Delete: %w", err)
	}

	if checkId != nil {

		queryDeleteWinner := fmt.Sprintf("UPDATE %s set is_winner = $1 WHERE id = $2", checks)

		_, err := tx.Exec(queryDeleteWinner, false, checkId)

		if err != nil {
			if txErr := tx.Rollback(); txErr != nil {
				return fmt.Errorf("repository.Delete: %w", txErr)
			}
			return fmt.Errorf("repository.Delete: %w", err)
		}

		queryDeleteMembers := fmt.Sprintf("DELETE FROM %s WHERE raffle_id = $1", members)

		_, err = tx.Exec(queryDeleteMembers, id)
		if err != nil {
			if txErr := tx.Rollback(); txErr != nil {
				return fmt.Errorf("repository.Delete: %w", txErr)
			}
			return fmt.Errorf("repository.Delete: %w", err)
		}

	}
	if txErr := tx.Commit(); txErr != nil {
		return fmt.Errorf("repository.Delete: %w", txErr)
	}
	return nil

}
