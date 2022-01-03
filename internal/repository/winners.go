package repository

import (
	"HundredToFive/internal/domain"
	"HundredToFive/pkg/logger"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"
)

type WinnerRepos struct {
	db *sqlx.DB
}

func NewWinnerRepos(db *sqlx.DB) *WinnerRepos {
	return &WinnerRepos{db: db}
}

func (w *WinnerRepos) CreateWinner(input domain.WinnerInput) error {

	var (
		inp             []*domain.CheckInfo
		raffleMembers   []*domain.Members
		raffleList      domain.Raffle
		err             error
		whereClause     string
		forRaffleValues []string
		setValues       string
	)

	if input.EndCheckDate == 0 {
		input.EndCheckDate = float64(time.Now().Unix())
	}
	if input.EndRegisteredDate == 0 {
		input.EndRegisteredDate = float64(time.Now().Unix())
	}

	if input.PartnerId != 0 {
		forRaffleValues = append(forRaffleValues, fmt.Sprintf("partner_id = %d", input.PartnerId))
	}

	forRaffleValues = append(forRaffleValues, fmt.Sprintf("check_amount <= %d", input.MoneyAmount))

	forRaffleValues = append(forRaffleValues, fmt.Sprintf("registered_at between to_timestamp(%f)::timestamp and to_timestamp(%f)::timestamp", input.StartRegisteredDate, input.EndRegisteredDate))

	forRaffleValues = append(forRaffleValues, fmt.Sprintf("check_date between to_timestamp(%f)::date and to_timestamp(%f)::date", input.StartCheckDate, input.EndCheckDate))
	forRaffleValues = append(forRaffleValues, fmt.Sprintf("is_winner = %s", false))
	whereClause = strings.Join(forRaffleValues, " AND ")

	if whereClause != "" {
		setValues = "WHERE " + whereClause
	}

	queryGetAllMembers := fmt.Sprintf(
		`SELECT 
				  	id  
				FROM 
					%s
				%s
				`, checks, setValues)

	err = w.db.Select(&inp, queryGetAllMembers)

	if err != nil {
		return fmt.Errorf("repository.CreateWinner:%w", err)
	}

	queryCheckExistWinner := fmt.Sprintf("SELECT check_id FROM %s WHERE id = $1", raffles)

	err = w.db.Get(&raffleList, queryCheckExistWinner, input.RaffleId)

	if err != nil {
		return fmt.Errorf("repository.CreateWinner:%w", err)
	}
	logger.Info(raffleList.CheckId)
	if raffleList.CheckId != nil {
		return fmt.Errorf("repository.CreateWinner: %w", domain.ErrWinnerAlreadyExistInRaffle)
	}

	tx := w.db.MustBegin()

	queryUpdateWinner := fmt.Sprintf("UPDATE %s SET check_id = $1,status = $2 WHERE id = $3", raffles)

	_, err = tx.Exec(queryUpdateWinner, input.CheckId, finished, input.RaffleId)

	if err != nil {

		if txErr := tx.Rollback(); txErr != nil {
			return fmt.Errorf("repository.CreateWinner: %w", txErr)
		}

		return fmt.Errorf("repository.CreateWinner: %w", err)
	}

	queryUpdateWinnerInCheck := fmt.Sprintf("UPDATE %s set is_winner = $1 WHERE id = $2", checks)

	_, err = tx.Exec(queryUpdateWinnerInCheck, true, input.CheckId)

	if err != nil {

		if txErr := tx.Rollback(); txErr != nil {
			return fmt.Errorf("repository.CreateWinner: %w", txErr)
		}

		return fmt.Errorf("repository.CreateWinner: %w", err)
	}

	for _, value := range inp {

		var member = domain.Members{
			CheckId:  value.Id,
			RaffleId: input.RaffleId,
		}

		raffleMembers = append(raffleMembers, &member)
	}

	queryInsertMembers := fmt.Sprintf("INSERT INTO %s(check_id,raffle_id) VALUES(:check_id,:raffle_id)", members)

	_, err = tx.NamedExec(queryInsertMembers, raffleMembers)

	if err != nil {

		if txErr := tx.Rollback(); txErr != nil {
			return fmt.Errorf("repository.CreateWinner: %w", txErr)
		}

		return fmt.Errorf("repository.CreateWinner: %w", err)
	}

	queryUpdateWinnerInRaffle := fmt.Sprintf("UPDATE %s set is_winner = $1 WHERE check_id = $2 AND raffle_id = $3", members)

	_, err = tx.Exec(queryUpdateWinnerInRaffle, true, input.CheckId, input.RaffleId)

	if err != nil {

		if txErr := tx.Rollback(); txErr != nil {
			return fmt.Errorf("repository.CreateWinner: %w", txErr)
		}

		return fmt.Errorf("repository.CreateWinner: %w", err)
	}

	txErr := tx.Commit()

	if txErr != nil {
		return fmt.Errorf("repository.CreateWinner: %w", txErr)
	}

	return nil

}

func (w *WinnerRepos) GetAll(page domain.Pagination, date int64) (*domain.GetAllWinnersCategoryResponse, error) {

	var (
		setValues string

		count int
	)
	logger.Info(date)
	if date != 0 {
		setValues = fmt.Sprintf("WHERE date(r.raffle_date) = to_timestamp(%d)::date", date)
	}

	queryCount := fmt.Sprintf(`SELECT COUNT(*) FROM %s r INNER JOIN %s ch on (r.check_id = ch.id and ch.is_winner) %s`, raffles, checks, setValues)

	err := w.db.QueryRowx(queryCount).Scan(&count)

	if err != nil {
		return nil, fmt.Errorf("repository.GetAll: %w", err)
	}

	offset, pagesCount := calculatePagination(&page, count)

	inp := make([]*domain.Winners, 0, page.Limit)

	query := fmt.Sprintf(
		`SELECT 
					r.id,
					us.phone_number,
					us.user_name,
					r.raffle_type,
					ch.is_winner,
					r.check_id,
					extract(epoch from r.raffle_date) "raffle_date"
				FROM 
					%s r
				inner join 
					%s ch
				on  
					r.check_id = ch.id
				inner join 
					%s us
				on 
					ch.user_id = us.id
				%s 
					ORDER BY r.id ASC LIMIT $1 OFFSET $2
			`, raffles, checks, usersTable, setValues)

	err = w.db.Select(&inp, query, page.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("repository.GetAll: %w", err)
	}

	pages := domain.PaginationPage{
		Page:  page.Page,
		Pages: pagesCount,
		Count: count,
	}

	winner := domain.GetAllWinnersCategoryResponse{
		Data:     inp,
		PageInfo: pages,
	}
	return &winner, nil
}

func (w *WinnerRepos) GetAllMembers(page domain.Pagination, id int) (*domain.GetAllWinnersCategoryResponse, error) {
	var (
		setValues string
		count     int
	)

	if id != 0 {
		setValues = fmt.Sprintf("WHERE mem.raffle_id = %d", id)
	}

	queryCount := fmt.Sprintf("SELECT COUNT(*) FROM %s mem %s", members, setValues)

	err := w.db.QueryRowx(queryCount).Scan(&count)

	if err != nil {
		return nil, fmt.Errorf("repository.GetAll: %w", err)
	}
	offset, pagesCount := calculatePagination(&page, count)

	inp := make([]*domain.Winners, 0, page.Limit)

	query := fmt.Sprintf(
		`SELECT 
					mem.id,
					mem.check_id,	
					us.phone_number,
					us.user_name,
					mem.is_winner,
					r.raffle_type,
					extract(epoch from r.raffle_date) "raffle_date"
				FROM 
					%s mem
				inner join
					%s r
				on 
					mem.raffle_id = r.id
				inner join 
					%s ch
				on  
					mem.check_id = ch.id
				inner join 
					%s us
				on 
					ch.user_id = us.id
				%s 
					ORDER BY mem.id ASC LIMIT $1 OFFSET $2
			`, members, raffles, checks, usersTable, setValues)

	err = w.db.Select(&inp, query, page.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("repository.GetAll: %w", err)
	}

	pages := domain.PaginationPage{
		Page:  page.Page,
		Pages: pagesCount,
		Count: count,
	}

	member := domain.GetAllWinnersCategoryResponse{
		Data:     inp,
		PageInfo: pages,
	}
	return &member, nil

}

func (w *WinnerRepos) GetAllDays(page domain.Pagination, month int) (*domain.GetAllDaysResponse, error) {
	var (
		setValues string
		count     int
	)

	if month != 0 {
		setValues = fmt.Sprintf("WHERE extract(month from raffle_date) = %d", month)
	}

	queryCount := fmt.Sprintf("SELECT count( distinct raffle_date) FROM %s %s", raffles, setValues)

	err := w.db.QueryRowx(queryCount).Scan(&count)
	logger.Info(count)

	if err != nil {
		return nil, fmt.Errorf("repository.GetAll: %w", err)
	}
	offset, pagesCount := calculatePagination(&page, count)

	inp := make([]*domain.DayInput, 0, page.Limit)

	query := fmt.Sprintf("SELECT extract(epoch from raffle_date) as created_at from %s %s group by raffle_date ORDER BY created_at DESC LIMIT $1 OFFSET $2", raffles, setValues)
	err = w.db.Select(&inp, query, page.Limit, offset)

	if err != nil {
		return nil, fmt.Errorf("repository.GetAllDays: %w", err)
	}

	pages := domain.PaginationPage{
		Page:  page.Page,
		Pages: pagesCount,
		Count: count,
	}

	date := domain.GetAllDaysResponse{
		Data:     inp,
		PageInfo: pages,
	}
	return &date, nil
}

func (w *WinnerRepos) GetAllMonths(page domain.Pagination) (*domain.GetAllDaysResponse, error) {
	var (
		count int
	)

	queryCount := fmt.Sprintf("SELECT COUNT(distinct extract(month from raffle_date)) FROM %s", raffles)
	err := w.db.QueryRowx(queryCount).Scan(&count)

	if err != nil {
		return nil, fmt.Errorf("repository.GetAllMonths: %w", err)
	}

	offset, pagesCount := calculatePagination(&page, count)

	inp := make([]*domain.DayInput, 0, page.Limit)

	query := fmt.Sprintf("select extract(month from raffle_date) as created_at from %s group by extract(month from raffle_date) ORDER BY created_at DESC LIMIT $1 OFFSET $2 ", raffles)
	err = w.db.Select(&inp, query, page.Limit, offset)

	if err != nil {
		return nil, fmt.Errorf("repository.GetAllMonths: %w", err)
	}
	pages := domain.PaginationPage{
		Page:  page.Page,
		Pages: pagesCount,
		Count: count,
	}

	date := domain.GetAllDaysResponse{
		Data:     inp,
		PageInfo: pages,
	}
	return &date, nil
}
