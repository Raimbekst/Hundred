package repository

import (
	"HundredToFive/internal/domain"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"math"
)

const (
	dayGame            = "Ежедневный розыгрыш"
	weeklyGame         = "Еженедельный розыгрыш"
	monthlyGame        = "Ежемесячный розыгрыш"
	finished           = "состоялся"
	notFinished        = "не состоялся"
	planned            = "запланирован"
	true               = "true"
	false              = "false"
	usersTable         = "users"
	sessionTable       = "sessions"
	cities             = "cities"
	partners           = "partners"
	banners            = "banners"
	checks             = "checks"
	checkImages        = "check_images"
	competition        = "competitions"
	members            = "members"
	raffles            = "raffles"
	userType           = "user"
	adminType          = "admin"
	faqs               = "faqs"
	descriptions       = "descriptions"
	websiteLinks       = "about"
	notifications      = "notifications"
	getters            = "getters"
	conditions         = "conditions"
	notificationTokens = "notification_users"
)

type UserAuth interface {
	VerifyExistenceUser(phone, email string) error

	CreateUser(user domain.User) (int, error)
	CreateAdmin(user domain.User) (int, error)

	SignIn(phone, password string) (*domain.User, error)
	SetSession(userId int, session domain.Session) error
	VerifyViaEmail(email string) (domain.User, error)
	ResetPassword(user domain.User) error
	GetAll(page domain.Pagination) (*domain.GetAllUsersCategoryResponse, error)
	GetByRefreshToken(refreshToken string) (domain.User, error)
	UserMe(id int) (*domain.UserList, error)

	BlockUser(id int) error
	UnBlockUser(id int) error
}

type City interface {
	Create(city domain.City) (int, error)
	GetAll(page domain.Pagination) (domain.GetAllCityCategoryResponse, error)
	GetById(id int) (domain.City, error)
	Update(id int, inp domain.City) error
	Delete(id int) error
}

type Partner interface {
	Create(partner domain.Partner) (int, error)
	GetAll(page domain.Pagination, status int) (*domain.GetAllPartnersCategoryResponse, error)
	GetById(id int) (domain.Partner, error)
	Update(id int, inp domain.UpdatePartner) ([]string, error)
	Delete(id int) ([]string, error)
}

type Banner interface {
	Create(banner domain.Banner) (int, error)
	GetAll(page domain.Pagination, status int, lang string) (*domain.GetAllBannersCategoryResponse, error)
	GetById(id int) (domain.Banner, error)
	Update(id int, inp domain.Banner) (string, error)
	Delete(id int) (string, error)
}

type Check interface {
	Create(info domain.CheckInfo) (int, error)
	GetAll(c *fiber.Ctx, page domain.Pagination, filter domain.FilterForCheck) (*domain.GetAllChecksCategoryResponse, error)
	GetById(ctx *fiber.Ctx, id int) (*domain.UserChecks, error)
}

type Winner interface {
	CreateWinner(input domain.WinnerInput) error
	GetAll(page domain.Pagination, id int) (*domain.GetAllWinnersCategoryResponse, error)
	GetAllMembers(page domain.Pagination, id int) (*domain.GetAllWinnersCategoryResponse, error)
	GetAllDays(page domain.Pagination, month int) (*domain.GetAllDaysResponse, error)
	GetAllMonths(page domain.Pagination) (*domain.GetAllDaysResponse, error)
}
type Raffle interface {
	Create(city domain.Raffle) (int, error)
	GetAll(page domain.Pagination, filter domain.FilterForRaffles) (*domain.GetAllRaffleCategoryResponse, error)
	GetById(id int) (domain.Raffle, error)
	Update(id int, inp domain.UpdateRaffle) error
	Delete(id int) error

	UpdateStatus(timeNow int64) error
}

type Faq interface {
	Create(faq domain.Faq) (int, error)
	GetAll(page domain.Pagination, lang string) (*domain.GetAllFaqsCategoryResponse, error)
	GetById(id int) (domain.Faq, error)
	Update(id int, inp domain.Faq) error
	Delete(id int) error

	CreateDesc(desc domain.Description) (int, error)
	GetAllDesc(page domain.Pagination, lang string) (*domain.GetAllDescCategoryResponse, error)
	GetDescById(id int) (domain.Description, error)
	UpdateDesc(id int, inp domain.Description) error
	DeleteDesc(id int) error
}

type Notification interface {
	Create(noty domain.Notification) (int, error)
	CreateForUser(noty domain.Notification) ([]string, int, error)

	GetAll(page domain.Pagination) (*domain.GetAllNotificationsResponse, error)

	GetById(id int) (*domain.Notification, error)
	Update(id int, inp domain.Notification) error
	Delete(id int) error

	StoreUsersToken(userId *int, token string) (int, error)

	GetAllRegistrationTokens() ([]string, error)

	GetNotificationByDate(time int64) ([]domain.Notification, error)
}

type About interface {
	Create(about domain.AboutUs) error

	GetAll() ([]*domain.AboutUs, error)
	Update(about domain.AboutUs) error
	Delete() error
}

type Condition interface {
	Create(con domain.Condition) (int, error)
	GetAll(page domain.Pagination, lang string) (*domain.GetAllConditionCategoryResponse, error)
	GetById(id int) (domain.Condition, error)
	Update(id int, inp domain.Condition) error
	Delete(id int) error
}

type Repository struct {
	UserAuth     UserAuth
	City         City
	Partner      Partner
	Banner       Banner
	Check        Check
	Winner       Winner
	Raffle       Raffle
	Faq          Faq
	About        About
	Notification Notification
	Condition    Condition
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		UserAuth:     NewUserAuthRepos(db),
		City:         NewCityRepos(db),
		Partner:      NewPartnerRepos(db),
		Banner:       NewBannerRepos(db),
		Check:        NewCheckRepository(db),
		Raffle:       NewRaffleRepos(db),
		Faq:          NewFaqsRepos(db),
		About:        NewAboutRepos(db),
		Winner:       NewWinnerRepos(db),
		Notification: NewNotificationRepos(db),
		Condition:    NewConditionRepos(db),
	}
}

func calculatePagination(page *domain.Pagination, count int) (int, int) {
	if page.Limit == 0 {
		page.Limit = count
	}

	if page.Page == 0 {
		page.Page = 1
	}

	pagesCount := 1.0

	if count != 0 {
		pagesCount = math.Ceil(float64(count) / float64(page.Limit))
		if page.Limit >= count {
			pagesCount = 1
		}
	}

	offset := (page.Page - 1) * page.Limit

	return offset, int(pagesCount)
}

func countPage(db *sqlx.DB, ss string) (int, error) {

	var count int
	queryCount := fmt.Sprintf("SELECT COUNT(*) FROM %s", ss)

	row := db.QueryRow(queryCount)
	err := row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("repository.countPage : %w", err)
	}
	return count, nil
}
