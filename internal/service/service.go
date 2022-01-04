package service

import (
	"HundredToFive/internal/config"
	"HundredToFive/internal/domain"
	"HundredToFive/internal/repository"
	"HundredToFive/pkg/auth"
	"HundredToFive/pkg/email"
	"HundredToFive/pkg/hash"
	"HundredToFive/pkg/phone"
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
	"time"
)

const (
	active     = "активный"
	blocked    = "заблокирован"
	planned    = "запланирован"
	finished   = "состоялся"
	unfinished = "не состоялся"
)

type SignUpInput struct {
	FirstName       string
	PhoneNumber     string
	Email           string
	Gender          string
	Age             int
	City            int
	Password        string
	ConfirmPassword string
	UserType        string
}
type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type VerificationEmailInput struct {
	Email  string
	Name   string
	Token  string
	Domain string
}

type UserAuth interface {
	VerifyExistenceUser(phone, email string) error

	UserSignUp(input SignUpInput) (int, error)
	AdminSignUp(input SignUpInput) (int, error)

	UserSignIn(user domain.User) (*Tokens, error)
	VerifyEmail(email, domain string) (string, error)
	ResetPassword(id int, newPassword, newPasswordConfirm string) error
	GetAll(page domain.Pagination) (*domain.GetAllUsersCategoryResponse, error)
	RefreshToken(token string) (*Tokens, error)
	UserMe(id int) (*domain.UserList, error)

	BlockUser(id int) error
	UnBlockUser(id int) error

	DownloadUsers(file *excelize.File) (*excelize.File, error)
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
	Update(id int, inp domain.Partner) error
	Delete(id int) error

	DownloadPartners(file *excelize.File, url string) (*excelize.File, error)
}

type Banner interface {
	Create(banner domain.Banner) (int, error)
	GetAll(page domain.Pagination, status int) (*domain.GetAllBannersCategoryResponse, error)
	GetById(id int) (domain.Banner, error)
	Update(id int, inp domain.Banner) error
	Delete(id int) error
}

type Emails interface {
	SendUserVerificationEmail(VerificationEmailInput) error
}
type Check interface {
	Create(info domain.CheckInfo) (int, error)
	GetAll(c *fiber.Ctx, page domain.Pagination, filter domain.FilterForCheck) (*domain.GetAllChecksCategoryResponse, error)
	GetById(c *fiber.Ctx, id int) (*domain.UserChecks, error)

	DownloadChecks(ctx *fiber.Ctx, file *excelize.File, filter domain.FilterForCheck) (*excelize.File, error)
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
	Update(id int, inp domain.Raffle) error
	Delete(id int) error
}

type Faq interface {
	Create(faq domain.Faq) (int, error)
	GetAll(page domain.Pagination) (*domain.GetAllFaqsCategoryResponse, error)
	GetById(id int) (domain.Faq, error)
	Update(id int, inp domain.Faq) error
	Delete(id int) error

	CreateDesc(desc domain.Description) (int, error)
	GetAllDesc(page domain.Pagination) (*domain.GetAllDescCategoryResponse, error)
	GetDescById(id int) (domain.Description, error)
	UpdateDesc(id int, inp domain.Description) error
	DeleteDesc(id int) error
}

type About interface {
	Create(about domain.AboutUs) error

	GetAll() ([]*domain.AboutUs, error)
	Update(about domain.AboutUs) error
	Delete() error
}

type Notification interface {
	Create(noty domain.Notification) (int, error)
	CreateForUser(noty domain.Notification) ([]string, int, error)

	GetAll(page domain.Pagination) (*domain.GetAllNotificationsResponse, error)

	GetById(id int) (*domain.Notification, error)
	Update(id int, inp domain.Notification) error
	Delete(id int) error

	StoreUsersToken(token string) (int, error)
}

type Condition interface {
	Create(con domain.Condition) (int, error)
	GetAll(page domain.Pagination) (*domain.GetAllConditionCategoryResponse, error)
	GetById(id int) (domain.Condition, error)
	Update(id int, inp domain.Condition) error
	Delete(id int) error
}

type Service struct {
	UserAuth
	City
	Partner
	Banner
	Emails
	Check
	Winner
	Raffle
	Faq
	About
	Notification
	Condition
}

type Deps struct {
	Repos           *repository.Repository
	Hashes          hash.PasswordHashes
	OtpPhone        phone.SecretGenerator
	Ctx             context.Context
	Redis           *redis.Client
	TokenManager    auth.TokenManager
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	EmailSender     email.Sender
	EmailConfig     config.EmailConfig
}

func NewService(deps Deps) *Service {
	emailsService := NewEmailService(deps.EmailSender, deps.EmailConfig)
	return &Service{
		Emails:       emailsService,
		UserAuth:     NewUserAuthService(deps.Repos.UserAuth, deps.Hashes, deps.OtpPhone, deps.Redis, deps.Ctx, deps.TokenManager, deps.AccessTokenTTL, deps.RefreshTokenTTL, emailsService),
		City:         NewCityService(deps.Repos.City),
		Partner:      NewPartnerService(deps.Repos.Partner),
		Banner:       NewBannerService(deps.Repos.Banner),
		Check:        NewCheckService(deps.Repos.Check),
		Winner:       NewWinnerService(deps.Repos.Winner),
		Raffle:       NewRaffleService(deps.Repos.Raffle),
		Faq:          NewFaqService(deps.Repos.Faq),
		About:        NewAboutService(deps.Repos.About),
		Notification: NewNotificationService(deps.Repos.Notification),
		Condition:    NewConditionService(deps.Repos.Condition),
	}
}
