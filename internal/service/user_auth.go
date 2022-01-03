package service

import (
	"HundredToFive/internal/domain"
	"HundredToFive/internal/repository"
	"HundredToFive/pkg/auth"
	"HundredToFive/pkg/hash"
	"HundredToFive/pkg/logger"
	"HundredToFive/pkg/phone"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/xuri/excelize/v2"
	"strconv"
	"time"
)

type UserAuthService struct {
	repos           repository.UserAuth
	hashes          hash.PasswordHashes
	otpPhone        phone.SecretGenerator
	redis           *redis.Client
	ctx             context.Context
	tokenManager    auth.TokenManager
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	emailService    Emails
}

func NewUserAuthService(
	repos repository.UserAuth,
	hashes hash.PasswordHashes,
	otpPhone phone.SecretGenerator,
	redis *redis.Client,
	ctx context.Context,
	tokenManager auth.TokenManager,
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
	emailService Emails) *UserAuthService {
	return &UserAuthService{
		repos:           repos,
		hashes:          hashes,
		otpPhone:        otpPhone,
		redis:           redis,
		ctx:             ctx,
		tokenManager:    tokenManager,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
		emailService:    emailService,
	}
}

func (u *UserAuthService) VerifyExistenceUser(phone, email string) error {
	return u.repos.VerifyExistenceUser(phone, email)
}

func (u *UserAuthService) UserSignUp(input SignUpInput) (int, error) {
	if input.Password != input.ConfirmPassword {
		return 0, fmt.Errorf("%w", domain.ErrPasswordNotMatch)
	}
	hashedPassword, err := u.hashes.Hash(input.Password)
	if err != nil {
		return 0, fmt.Errorf("service.UserSignUp: %w", err)
	}

	user := domain.User{
		Name:        input.FirstName,
		PhoneNumber: input.PhoneNumber,
		Email:       input.Email,
		Gender:      input.Gender,
		Age:         input.Age,
		CityId:      input.City,
		Password:    hashedPassword,
		UserType:    input.UserType,
	}

	id, err := u.repos.CreateUser(user)
	if err != nil {
		return 0, fmt.Errorf("service.UserSignUp: %w", err)
	}

	return id, nil
}

func (u *UserAuthService) AdminSignUp(input SignUpInput) (int, error) {
	if input.Password != input.ConfirmPassword {
		return 0, fmt.Errorf("%w", domain.ErrPasswordNotMatch)
	}
	hashedPassword, err := u.hashes.Hash(input.Password)
	if err != nil {
		return 0, fmt.Errorf("service.UserSignUp: %w", err)
	}

	user := domain.User{
		Name:        input.FirstName,
		PhoneNumber: input.PhoneNumber,
		Email:       input.Email,
		Gender:      input.Gender,
		Age:         input.Age,
		Password:    hashedPassword,
		UserType:    input.UserType,
	}

	id, err := u.repos.CreateAdmin(user)
	if err != nil {
		return 0, fmt.Errorf("service.UserSignUp: %w", err)
	}

	return id, nil
}

func (u *UserAuthService) VerifyEmail(email, domain string) (string, error) {

	users, err := u.repos.VerifyViaEmail(email)
	if err != nil {
		return "", fmt.Errorf("service.VerifyEmail: %w", err)
	}

	token, err := u.tokenManager.NewJWT(strconv.Itoa(users.Id), users.UserType, u.accessTokenTTL)
	if err != nil {
		return "", fmt.Errorf("service.VerifyEmail:%w", err)
	}

	err = u.emailService.SendUserVerificationEmail(VerificationEmailInput{
		Email:  email,
		Name:   users.Name,
		Token:  token,
		Domain: domain,
	})

	if err != nil {
		return "", fmt.Errorf("service.VerifyEmail:%w", err)
	}
	return token, nil
}

func (u *UserAuthService) ResetPassword(id int, newPassword, newPasswordConfirm string) error {
	if newPassword != newPasswordConfirm {
		return fmt.Errorf("%w", domain.ErrPasswordNotMatch)
	}
	hashedPassword, err := u.hashes.Hash(newPassword)
	if err != nil {
		return fmt.Errorf("service.ResetPassword: %w", err)
	}

	user := domain.User{
		Id:       id,
		Password: hashedPassword,
	}
	return u.repos.ResetPassword(user)

}

func (u *UserAuthService) UserSignIn(user domain.User) (*Tokens, error) {
	hashedPassword, err := u.hashes.Hash(user.Password)
	if err != nil {
		return nil, fmt.Errorf("service.UserSignIn: %w", err)
	}
	logger.Info(hashedPassword)
	input, err := u.repos.SignIn(user.PhoneNumber, hashedPassword)
	if err != nil {
		return nil, fmt.Errorf("service.UserSignIn: %w", err)
	}
	return u.createSession(input.Id, input.UserType)
}

func (u *UserAuthService) UserMe(id int) (*domain.UserList, error) {
	return u.repos.UserMe(id)
}

func (u *UserAuthService) RefreshToken(refreshToken string) (*Tokens, error) {
	user, err := u.repos.GetByRefreshToken(refreshToken)

	if err != nil {
		return nil, fmt.Errorf("service.UserSignIn: %w", err)
	}
	res, err := u.createSession(user.Id, user.UserType)

	if err != nil {
		return nil, fmt.Errorf("service.RefreshToken: %w", err)
	}

	return res, nil
}

func (u *UserAuthService) createSession(userId int, userType string) (*Tokens, error) {
	var (
		res Tokens
		err error
	)

	res.AccessToken, err = u.tokenManager.NewJWT(strconv.Itoa(userId), userType, u.accessTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("service.createSession.NewJWT: %w", err)

	}

	res.RefreshToken, err = u.tokenManager.NewRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("service.createSession.NewRefreshToken: %w", err)
	}

	session := domain.Session{
		RefreshToken: res.RefreshToken,
		ExpiresAt:    time.Now().Add(u.refreshTokenTTL),
	}

	err = u.repos.SetSession(userId, session)

	if err != nil {
		return nil, fmt.Errorf("service.createSession: %w", err)
	}

	return &res, nil
}

func (u *UserAuthService) GetAll(page domain.Pagination) (*domain.GetAllUsersCategoryResponse, error) {
	return u.repos.GetAll(page)
}

func (u *UserAuthService) BlockUser(id int) error {
	return u.repos.BlockUser(id)
}

func (u *UserAuthService) UnBlockUser(id int) error {
	return u.repos.UnBlockUser(id)
}

func (u *UserAuthService) DownloadUsers(file *excelize.File) (*excelize.File, error) {
	page := domain.Pagination{}
	list, err := u.repos.GetAll(page)

	if err != nil {
		return nil, fmt.Errorf("service.DownloadUsers: %w", err)
	}
	id := 2

	for _, value := range list.Data {
		status := active
		if value.IsBlocked {

			status = blocked
		}
		unixTime := time.Unix(int64(value.RegisteredAt), 0)
		file.SetCellValue("Sheet1", "A"+strconv.Itoa(id), value.Id)
		file.SetCellValue("Sheet1", "B"+strconv.Itoa(id), value.Name)
		file.SetCellValue("Sheet1", "C"+strconv.Itoa(id), value.PhoneNumber)
		file.SetCellValue("Sheet1", "D"+strconv.Itoa(id), value.Email)
		file.SetCellValue("Sheet1", "E"+strconv.Itoa(id), value.Gender)
		file.SetCellValue("Sheet1", "F"+strconv.Itoa(id), value.Age)
		file.SetCellValue("Sheet1", "G"+strconv.Itoa(id), value.City)
		file.SetCellValue("Sheet1", "H"+strconv.Itoa(id), unixTime)
		file.SetCellValue("Sheet1", "I"+strconv.Itoa(id), status)
		id = id + 1
	}

	return file, nil

}
