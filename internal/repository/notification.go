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
		userTokens  []domain.NotificationToken
		tokens      []string
	)

	for _, value := range noty.Ids {
		userTokens = nil

		queryGetTokens := fmt.Sprintf("SELECT registration_token FROM %s WHERE user_id = $1 ", notificationTokens)

		err := n.db.Select(&userTokens, queryGetTokens, value)

		if err != nil {
			continue
		}

		for _, val := range userTokens {
			tokens = append(tokens, val.RegistrationToken)
		}

	}

	tx := n.db.MustBegin()

	query := fmt.Sprintf(
		`INSERT INTO 
						%s
					(title,text,link,noty_date,noty_time,status,noty_getters)
						VALUES
					($1,$2,$3,to_timestamp($4) at time zone 'GMT',$5,$6,$7) RETURNING id`, notifications)

	err := tx.QueryRowx(query, noty.Title, noty.Text, noty.Link, noty.Date, noty.Time, noty.Status, noty.Getters).Scan(&id)

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
						n.id,n.title,n.text,
						n.link,n.noty_time,
						extract(epoch from n.noty_date::timestamp at time zone 'GMT') "noty_date",
						n.status,n.noty_getters
					FROM %s n
					ORDER BY id ASC LIMIT $1 OFFSET $2`, notifications)

	err = n.db.Select(&inp, query, page.Limit, offset)

	if err != nil {
		return nil, fmt.Errorf("repository.GetAll: %w", err)
	}

	var gettersNotification domain.GetterList

	for _, value := range inp {

		queryGetAllUsers := fmt.Sprintf("SELECT * FROM %s WHERE notification_id = $1", getters)

		rows, err := n.db.Queryx(queryGetAllUsers, value.Id)
		if err != nil {
			return nil, fmt.Errorf("repository.GetAll: %w", err)
		}
		for rows.Next() {
			err = rows.StructScan(&gettersNotification)
			if err != nil {
				return nil, fmt.Errorf("repository.GetAll: %w", err)
			}
			value.Users = append(value.Users, gettersNotification)
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
						n.id,n.title,n.text,
						n.link,n.noty_time,
						extract(epoch from n.noty_date::timestamp at time zone 'GMT') "noty_date",
						n.status,n.noty_getters
					FROM %s n
					FULL OUTER JOIN %s p 
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

	var gettersNotification domain.GetterList
	for rows.Next() {
		err := rows.StructScan(&gettersNotification)
		if err != nil {
			return nil, fmt.Errorf("repository.GetById: %w", err)
		}
		noty.Users = append(noty.Users, gettersNotification)
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

	if inp.Time != 0 {
		setValues = append(setValues, fmt.Sprintf("noty_time = $%d", argId))
		args = append(args, inp.Time)
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

	if inp.PartnerId != nil {
		setValues = append(setValues, fmt.Sprintf("partner_id=$%d", argId))
		args = append(args, inp.PartnerId)
		argId++
	}
	if inp.Status != 0 {
		setValues = append(setValues, fmt.Sprintf("status = $%d", argId))
		args = append(args, inp.Status)
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

func (n *NotificationRepos) StoreUsersToken(userId *int, token string) (int, error) {
	var (
		id     int
		tokens domain.NotificationToken
	)

	queryCheck := fmt.Sprintf("SELECT id FROM %s WHERE registration_token = $1", notificationTokens)

	err := n.db.Get(&tokens, queryCheck, token)

	if err == nil {

		queryUpdate := fmt.Sprintf("UPDATE %s SET user_id = $1 WHERE registration_token = $2 RETURNING id", notificationTokens)

		err = n.db.QueryRowx(queryUpdate, userId, token).Scan(&id)
		if err != nil {
			return 0, fmt.Errorf("repository.StoreUsersToken: %w", err)
		}

	} else {

		query := fmt.Sprintf("INSERT INTO %s(user_id,registration_token) VALUES($1, $2) RETURNING id", notificationTokens)

		err = n.db.QueryRowx(query, userId, token).Scan(&id)
		if err != nil {
			return 0, fmt.Errorf("repository.StoreUsersToken: %w", err)
		}
	}
	return id, nil
}

func (n *NotificationRepos) GetAllRegistrationTokens() ([]string, error) {
	var (
		tokens    []domain.NotificationToken
		tokenList []string
	)

	query := fmt.Sprintf("SELECT * FROM %s", notificationTokens)

	err := n.db.Select(&tokens, query)

	if err != nil {
		return nil, fmt.Errorf("repository.GetAllRegistrationTokens: %w", err)
	}

	for _, value := range tokens {
		tokenList = append(tokenList, value.RegistrationToken)
	}

	return tokenList, nil
}

func (n *NotificationRepos) GetNotificationByDate(time int64) ([]domain.Notification, error) {
	var input []domain.Notification

	query := fmt.Sprintf(
		`SELECT 
					id,
					title,
					text,
					link,
					extract(epoch from noty_date::timestamp at time zone 'GMT') "noty_date"
				   FROM 
				%s 
				WHERE noty_date = to_timestamp($1) at time zone 'GMT'`, notifications)

	err := n.db.Select(&input, query, time)

	if err != nil {
		return nil, fmt.Errorf("repository.GetNotificationByDate: %w", err)
	}
	return input, nil
}
