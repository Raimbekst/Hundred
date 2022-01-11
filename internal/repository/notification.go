package repository

import (
	"HundredToFive/internal/domain"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
)

type NotificationRepos struct {
	db *sqlx.DB
}

func NewNotificationRepos(db *sqlx.DB) *NotificationRepos {
	return &NotificationRepos{db: db}
}

func (n *NotificationRepos) Create(noty domain.Notification) (int, error) {

	var id int

	query := fmt.Sprintf(
		`INSERT INTO 
						%s
					(title,text,partner_id,link,reference,noty_date,noty_time,status,noty_getters)	
						VALUES
					($1,$2,$3,$4,$5,to_timestamp($6) at time zone 'GMT',$7,$8,$9) RETURNING id`, notifications)

	err := n.db.QueryRowx(query, noty.Title, noty.Text, noty.PartnerId, noty.Link, noty.Reference, noty.Date, noty.Time, noty.Status, noty.Getters).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("repository.Create: %w", err)
	}

	return id, nil
}

func (n *NotificationRepos) CreateForUser(noty domain.Notification) ([]string, int, error) {

	var (
		id          int
		getterUsers []domain.GetterList
		userTokens  domain.NotificationToken
		tokens      []string
	)

	for _, value := range noty.Ids {
		queryGetTokens := fmt.Sprintf("SELECT registration_token FROM %s WHERE id = $1 ", notificationTokens)

		err := n.db.Get(&userTokens, queryGetTokens, value)

		if err != nil {
			continue
		}

		tokens = append(tokens, userTokens.RegistrationToken)
	}

	tx := n.db.MustBegin()

	query := fmt.Sprintf(
		`INSERT INTO 
						%s
					(title,text,link,noty_date,status,noty_getters)
						VALUES
					($1,$2,$3,to_timestamp($4) timestamp at time zone 'GMT',$5,$6) RETURNING id`, notifications)

	err := tx.QueryRowx(query, noty.Title, noty.Text, noty.Link, noty.Date, noty.Status, noty.Getters).Scan(&id)

	if err != nil {
		if txErr := tx.Rollback(); txErr != nil {
			return nil, 0, fmt.Errorf("repository.CreateForUser: %w", txErr)
		}
		return nil, 0, fmt.Errorf("repository.CreateForUser: %w", err)
	}

	for _, value := range noty.Ids {

		var getter = domain.GetterList{
			NotificationId: id,
			UserId:         value,
		}

		getterUsers = append(getterUsers, getter)
	}

	queryInsertUser := fmt.Sprintf("INSERT INTO %s(notification_id,user_id) VALUES(:notification_id,:user_id)", getters)

	rows, err := tx.NamedExec(queryInsertUser, getterUsers)

	if err != nil {
		if txErr := tx.Rollback(); txErr != nil {
			return nil, 0, fmt.Errorf("repository.CreateForUser: %w", txErr)
		}
		return nil, 0, fmt.Errorf("repository.CreateForUser: %w", err)
	}
	affected, err := rows.RowsAffected()
	if err != nil {
		if txErr := tx.Rollback(); txErr != nil {
			return nil, 0, fmt.Errorf("repository.CreateForUser: %w", txErr)
		}
		return nil, 0, fmt.Errorf("repository.CreateForUser: %w", err)
	}

	if affected == 0 {
		if txErr := tx.Rollback(); txErr != nil {
			return nil, 0, fmt.Errorf("repository.CreateForUser: %w", txErr)
		}
		return nil, 0, fmt.Errorf("repository.CreateForUser: %w", sql.ErrNoRows)
	}
	if txErr := tx.Commit(); txErr != nil {
		return nil, 0, fmt.Errorf("repository.CreateForUser: %w", txErr)
	}

	return tokens, id, nil

}

func (n *NotificationRepos) GetAll(page domain.Pagination) (*domain.GetAllNotificationsResponse, error) {

	count, err := countPage(n.db, notifications)
	if err != nil {
		return nil, fmt.Errorf("repository.GetAll: %w", err)
	}

	offset, pagesCount := calculatePagination(&page, count)

	inp := make([]*domain.Notification, 0, page.Limit)

	query := fmt.Sprintf(
		`SELECT 
						n.id,n.title,n.text,p.logo,
						n.link,n.reference,
						extract(epoch from n.noty_date::timestamp at time zone 'GMT') "noty_date",
						n.noty_time,n.status,n.noty_getters
					FROM %s n
					INNER JOIN %s p on n.partner_id = p.id  
					ORDER BY id ASC LIMIT $1 OFFSET $2`, notifications, partners)

	err = n.db.Select(&inp, query, page.Limit, offset)

	if err != nil {
		return nil, fmt.Errorf("repository.GetAll: %w", err)
	}

	for _, value := range inp {

		queryGetAllUsers := fmt.Sprintf("SELECT * FROM %s WHERE notification_id = $1", getters)

		rows, err := n.db.Queryx(queryGetAllUsers, value.Id)
		if err != nil {
			return nil, fmt.Errorf("repository.GetAll: %w", err)
		}

		for rows.Next() {
			err := rows.StructScan(&value.Users)
			if err != nil {
				return nil, fmt.Errorf("repository.GetAll: %w", err)
			}
		}

	}

	pages := domain.PaginationPage{
		Page:  page.Page,
		Pages: pagesCount,
		Count: count,
	}

	noty := domain.GetAllNotificationsResponse{
		Data:     inp,
		PageInfo: pages,
	}
	return &noty, nil
}

func (n *NotificationRepos) GetById(id int) (*domain.Notification, error) {
	var noty domain.Notification

	query := fmt.Sprintf(
		`SELECT 
						n.id,n.title,n.text,p.logo,
						n.link,n.reference,
						extract(epoch from n.noty_date::timestamp at time zone 'GMT') "noty_date",
						n.noty_time,n.status,n.noty_getters
					FROM %s n
					INNER JOIN %s p 
					on n.partner_id = p.id WHERE n.id = $1`, notifications, partners)

	err := n.db.Get(&noty, query, id)
	if err != nil {
		return nil, fmt.Errorf("repository.GetById: %w", domain.ErrNotFound)
	}

	queryGetAllUsers := fmt.Sprintf("SELECT * FROM %s WHERE notification_id = $1", getters)

	rows, err := n.db.Queryx(queryGetAllUsers, noty.Id)
	if err != nil {
		return nil, fmt.Errorf("repository.GetById: %w", err)
	}
	for rows.Next() {
		err := rows.StructScan(&noty.Users)
		if err != nil {
			return nil, fmt.Errorf("repository.GetById: %w", err)
		}
	}

	return &noty, nil

}

func (n *NotificationRepos) Update(id int, inp domain.Notification) error {
	var input domain.Notification

	queryCheck := fmt.Sprintf("SELECT status FROM %s WHERE id = $1", notifications)

	err := n.db.Get(&input, queryCheck, id)

	if err != nil {
		return fmt.Errorf("repository.Update: %w", domain.ErrNotFound)
	}

	if input.Status != 1 {
		return fmt.Errorf("repository.Update: %w", domain.ErrUpdateNotification)
	}

	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if inp.Date != 0 {
		setValues = append(setValues, fmt.Sprintf(" noty_date = to_timestamp($%d) at time zone 'GMT' ", argId))
		args = append(args, inp.Date)
		argId++
	}

	if inp.Title != "" {
		setValues = append(setValues, fmt.Sprintf("title = $%d", argId))
		args = append(args, inp.Title)
		argId++
	}

	if inp.Text != "" {
		setValues = append(setValues, fmt.Sprintf("text=$%d", argId))
		args = append(args, inp.Text)
		argId++
	}

	if inp.Link != "" {
		setValues = append(setValues, fmt.Sprintf("link=$%d", argId))
		args = append(args, inp.Link)
		argId++
	}

	if inp.Reference != "" {
		setValues = append(setValues, fmt.Sprintf("reference=$%d", argId))
		args = append(args, inp.Reference)
		argId++
	}

	if inp.Time != 0 {
		setValues = append(setValues, fmt.Sprintf("noty_time=$%d", argId))
		args = append(args, inp.Time)
		argId++
	}

	if inp.PartnerId != 0 {
		setValues = append(setValues, fmt.Sprintf("partner_id=$%d", argId))
		args = append(args, inp.PartnerId)
		argId++
	}
	fmt.Println("sds")

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`UPDATE %s SET %s WHERE id = $%d `, notifications, setQuery, argId)
	args = append(args, id)

	result, err := n.db.Exec(query, args...)

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("repository.Update: %w", err)
	}

	if affected == 0 {
		return fmt.Errorf("repository.Update: %w", domain.ErrNotFound)

	}

	return err

}

func (n *NotificationRepos) Delete(id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", notifications)
	result, err := n.db.Exec(query, id)
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("repository.Delete: %w", err)
	}

	if affected == 0 {
		return fmt.Errorf("repository.Delete: %w", domain.ErrNotFound)
	}

	return nil
}

func (n *NotificationRepos) StoreUsersToken(token string) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s(registration_token) VALUES($1) RETURNING id", notificationTokens)

	err := n.db.QueryRowx(query, token).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("repository.StoreUsersToken: %w", err)
	}
	return id, nil
}
