package repository

import (
	"HundredToFive/internal/domain"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"reflect"
	"strings"
)

type UserAuthRepos struct {
	db *sqlx.DB
}

func NewUserAuthRepos(db *sqlx.DB) *UserAuthRepos {
	return &UserAuthRepos{db: db}
}

func (u *UserAuthRepos) VerifyExistenceUser(phone, email string) error {
	var user domain.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE phone_number = $1 or email = $2", usersTable)
	err := u.db.Get(&user, query, phone, email)
	if err != nil {
		return fmt.Errorf("repository.VerifyExistenceUser: %w", err)
	}
	return nil
}

func (u *UserAuthRepos) CreateUser(user domain.User) (int, error) {
	var city domain.City
	queryCity := fmt.Sprintf("SELECT name FROM %s WHERE id = $1", cities)
	err := u.db.Get(&city, queryCity, user.CityId)
	if err != nil {
		return 0, fmt.Errorf("repository.CreateUser: %w", domain.ErrCityNotFound)
	}
	tx := u.db.MustBegin()

	var id int
	query := fmt.Sprintf(
		`INSERT INTO 
					%s
				(user_name,phone_number,email,age,city,gender,password,user_type)
					VALUES($1,$2,$3,$4,$5,$6,$7,$8)
				RETURNING id`,
		usersTable)

	err = tx.QueryRowx(query, user.Name, user.PhoneNumber, user.Email, user.Age, city.Name, user.Gender, user.Password, user.UserType).Scan(&id)
	if err != nil {

		if txErr := tx.Rollback(); txErr != nil {
			return 0, fmt.Errorf("repository.CreateUser: %w", txErr)
		}
		return 0, fmt.Errorf("repository.CreateUser: %w", err)
	}

	querySession := fmt.Sprintf("INSERT INTO %s(user_id) VALUES($1)", sessionTable)
	_, err = tx.Exec(querySession, id)

	if err != nil {
		if txErr := tx.Rollback(); txErr != nil {
			return 0, fmt.Errorf("repository.CreateUser: %w", txErr)
		}
		return 0, fmt.Errorf("repository.CreateUser: %w", err)
	}

	if txErr := tx.Commit(); txErr != nil {
		return 0, fmt.Errorf("repository.CreateUser: %w", txErr)
	}
	return id, nil
}

func (u *UserAuthRepos) CreateAdmin(user domain.User) (int, error) {
	tx := u.db.MustBegin()

	var id int
	query := fmt.Sprintf(
		`INSERT INTO 
					%s
				(user_name,phone_number,email,age,gender,password,user_type)
					VALUES($1,$2,$3,$4,$5,$6,$7)
				RETURNING id`,
		usersTable)

	err := tx.QueryRowx(query, user.Name, user.PhoneNumber, user.Email, user.Age, user.Gender, user.Password, user.UserType).Scan(&id)
	if err != nil {

		if txErr := tx.Rollback(); txErr != nil {
			return 0, fmt.Errorf("repository.CreateAdmin: %w", txErr)
		}
		return 0, fmt.Errorf("repository.CreateAdmin: %w", err)
	}

	querySession := fmt.Sprintf("INSERT INTO %s(user_id) VALUES($1)", sessionTable)
	_, err = tx.Exec(querySession, id)

	if err != nil {
		if txErr := tx.Rollback(); txErr != nil {
			return 0, fmt.Errorf("repository.CreateAdmin: %w", txErr)
		}
		return 0, fmt.Errorf("repository.CreateAdmin: %w", err)
	}

	if txErr := tx.Commit(); txErr != nil {
		return 0, fmt.Errorf("repository.CreateAdmin: %w", txErr)
	}
	return id, nil

}

func (u *UserAuthRepos) SignIn(phone, password string) (*domain.User, error) {
	var input domain.User
	query := fmt.Sprintf("SELECT id,user_type,is_blocked FROM %s WHERE phone_number = $1 AND password = $2 ", usersTable)
	err := u.db.Get(&input, query, phone, password)
	if err != nil {
		return nil, fmt.Errorf("repository.SignIn: %w", domain.ErrUserNotExist)
	}
	if input.IsBlocked {
		return nil, fmt.Errorf("repository.SignIn: %w", domain.ErrUserBlocked)
	}
	return &input, nil
}

func (u *UserAuthRepos) UserMe(id int) (*domain.UserList, error) {
	var inp domain.UserList
	query := fmt.Sprintf("SELECT id,user_name,phone_number,email,age,gender,city FROM %s WHERE id = $1", usersTable)

	err := u.db.Get(&inp, query, id)

	if err != nil {
		return nil, fmt.Errorf("repository.UserMe: %w", err)
	}
	return &inp, nil
}

func (u *UserAuthRepos) GetByRefreshToken(refreshToken string) (domain.User, error) {
	var input domain.User

	query := fmt.Sprintf("SELECT u.id,u.user_type FROM %s u INNER JOIN %s s on u.id=s.user_id WHERE s.refresh_token=$1", usersTable, sessionTable)

	err := u.db.Get(&input, query, refreshToken)

	if err != nil {
		return domain.User{}, fmt.Errorf("repository.GetByRefreshToken: %w", err)
	}
	return input, nil
}

func (u *UserAuthRepos) VerifyViaEmail(email string) (domain.User, error) {

	var input domain.User
	query := fmt.Sprintf("SELECT id,user_name,user_type FROM %s WHERE email = $1 ", usersTable)

	err := u.db.Get(&input, query, email)

	if err != nil {
		return domain.User{}, fmt.Errorf("repository.VerifyViaEmail: %w", domain.ErrUserNotExist)
	}
	return input, nil
}

func (u *UserAuthRepos) SetSession(userId int, session domain.Session) error {
	setValues := make([]string, 0, reflect.TypeOf(domain.Session{}).NumField())

	if session.RefreshToken != "" {
		setValues = append(setValues, fmt.Sprintf("refresh_token=:refresh_token"))
	}
	if !session.ExpiresAt.IsZero() {
		setValues = append(setValues, fmt.Sprintf("expires_at=:expires_at"))
	}

	setQuery := strings.Join(setValues, ", ")

	if setQuery == "" {
		return fmt.Errorf("repository.SetSession: %v", "empty body")
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE user_id=%d", sessionTable, setQuery, userId)
	result, err := u.db.NamedExec(query, session)
	if err != nil {
		return fmt.Errorf("repository.SetSession: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("repository.SetSession: %w", err)
	}

	if affected == 0 {
		return fmt.Errorf("repository.Update: %w", sql.ErrNoRows)
	}

	return nil
}

func (u *UserAuthRepos) ResetPassword(user domain.User) error {
	query := fmt.Sprintf("UPDATE %s SET password = $1 WHERE id = $2", usersTable)
	rows, err := u.db.Exec(query, user.Password, user.Id)

	if err != nil {
		return fmt.Errorf("repository.ResetPassword: %w", err)
	}

	affected, err := rows.RowsAffected()

	if err != nil {
		return fmt.Errorf("repository.ResetPassword: %w", err)
	}

	if affected == 0 {
		return fmt.Errorf("repository.ResetPassword: %w", sql.ErrNoRows)
	}

	return nil
}

func (u *UserAuthRepos) GetAll(page domain.Pagination) (*domain.GetAllUsersCategoryResponse, error) {
	var (
		count int
	)

	queryCount := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE user_type = $1", usersTable)

	err := u.db.QueryRowx(queryCount, userType).Scan(&count)

	if err != nil {
		return nil, fmt.Errorf("repository.GetAll: %w", err)
	}

	offset, pagesCount := calculatePagination(&page, count)

	inp := make([]*domain.UserList, 0, page.Limit)

	query := fmt.Sprintf(
		`SELECT 
					id,user_name,phone_number,email,age,gender,city,is_blocked,extract(epoch from registered_at) "registered_at"
				FROM 
					%s WHERE user_type = $1 ORDER BY id ASC LIMIT $2 OFFSET $3`, usersTable)
	err = u.db.Select(&inp, query, userType, page.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("repository.GetAll: %w", err)
	}

	pages := domain.PaginationPage{
		Page:  page.Page,
		Pages: pagesCount,
		Count: count,
	}

	users := domain.GetAllUsersCategoryResponse{
		Data:     inp,
		PageInfo: pages,
	}
	return &users, nil
}

func (u *UserAuthRepos) BlockUser(id int) error {
	query := fmt.Sprintf("update %s set is_blocked = $1 where id=$2 and user_type = $3", usersTable)
	result, err := u.db.Exec(query, true, id, userType)

	if err != nil {
		return fmt.Errorf("repository.BlcokUser: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("repository.BlcokUser: %w", err)
	}

	if affected == 0 {
		return fmt.Errorf("repository.BlcokUser: %w", domain.ErrNotFound)
	}
	return nil

}

func (u *UserAuthRepos) UnBlockUser(id int) error {
	query := fmt.Sprintf("update %s set is_blocked = $1 where id=$2 and user_type = $3", usersTable)
	result, err := u.db.Exec(query, false, id, userType)

	if err != nil {
		return fmt.Errorf("repository.BlcokUser: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("repository.BlcokUser: %w", err)
	}

	if affected == 0 {
		return fmt.Errorf("repository.BlcokUser: %w", domain.ErrNotFound)
	}
	return nil
}
