package v1

import (
	"HundredToFive/internal/domain"
	"HundredToFive/internal/service"
	"HundredToFive/pkg/excel"
	"HundredToFive/pkg/validation/validationStructs"
	"bytes"
	"errors"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"strconv"
	"time"
)

func (h *Handler) initUserRoutes(api fiber.Router) {
	api.Post("/admin/sign-up", h.adminSignUp)

	auth := api.Group("/auth/")
	{
		auth.Get("users", h.getAllUsers)
		auth.Get("download/users", h.downloadUsers)
		auth.Post("sign-up", h.userSignUp)
		auth.Post("sign-in", h.userSignIn)
		auth.Post("verify", h.verifyEmail)
		auth.Post("reset-password", h.resetPassword)
		auth.Post("refresh", h.userRefresh)

		user := auth.Group("", jwtware.New(jwtware.Config{
			SigningKey: []byte(h.signingKey),
		}))
		{
			user.Get("user/me", h.userMe)
			user.Put("block-user/:id", h.blockUser)
			user.Put("unblock-user/:id", h.unblockUser)
		}
	}
}

type AdminSignUpInput struct {
	FirstName       string `json:"name" validate:"required"`
	PhoneNumber     string `json:"phone_number" db:"phone_number" validate:"required,e164"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8,max=64"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
}

// @Tags auth
// @Description create Admin
// @ModuleID adminSignUp
// @Accept json
// @Produce  json
// @Param data body AdminSignUpInput true "admin sign-up"
// @Success 201 {object} idResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admin/sign-up [post]
func (h *Handler) adminSignUp(c *fiber.Ctx) error {
	var input AdminSignUpInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	ok, errs := validationStructs.ValidateStruct(input)

	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(errs)
	}
	err := h.services.UserAuth.VerifyExistenceUser(input.PhoneNumber, input.Email)

	if err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: domain.ErrUserAlreadyExist.Error()})
	}

	gender := map[int]string{1: "male", 2: "female"}

	id, err := h.services.UserAuth.AdminSignUp(service.SignUpInput{
		FirstName:       input.FirstName,
		PhoneNumber:     input.PhoneNumber,
		Email:           input.Email,
		Gender:          gender[1],
		Password:        input.Password,
		ConfirmPassword: input.ConfirmPassword,
		UserType:        admin,
	})

	if err != nil {
		if errors.Is(err, domain.ErrPasswordNotMatch) || errors.Is(err, domain.ErrCityNotFound) || errors.Is(err, domain.ErrCityNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(idResponse{ID: id})
}

type UserSignUpInput struct {
	FirstName       string `json:"name" validate:"required"`
	PhoneNumber     string `json:"phone_number" db:"phone_number" validate:"required,e164"`
	Email           string `json:"email" validate:"required,email"`
	Gender          int    `json:"gender" form:"gender" validate:"required" minimum:"1" maximum:"2"`
	Age             int    `json:"age" validate:"required,gt=15,lt=91" examples:"24"`
	City            int    `json:"city" validate:"required"`
	Password        string `json:"password" validate:"required,min=8,max=64"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
}

// @Tags auth
// @Description create member account
// @ModuleID userSignUp
// @Accept json
// @Produce  json
// @Param data body UserSignUpInput true "user sign-up"
// @Success 201 {object} idResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /auth/sign-up [post]
func (h *Handler) userSignUp(c *fiber.Ctx) error {
	var input UserSignUpInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	ok, errs := validationStructs.ValidateStruct(input)

	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(errs)
	}

	err := h.services.UserAuth.VerifyExistenceUser(input.PhoneNumber, input.Email)

	if err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: domain.ErrUserAlreadyExist.Error()})
	}

	gender := map[int]string{1: "male", 2: "female"}

	id, err := h.services.UserAuth.UserSignUp(service.SignUpInput{
		FirstName:       input.FirstName,
		PhoneNumber:     input.PhoneNumber,
		Email:           input.Email,
		Gender:          gender[input.Gender],
		Age:             input.Age,
		City:            input.City,
		Password:        input.Password,
		ConfirmPassword: input.ConfirmPassword,
		UserType:        user,
	})

	if err != nil {
		if errors.Is(err, domain.ErrPasswordNotMatch) || errors.Is(err, domain.ErrCityNotFound) || errors.Is(err, domain.ErrCityNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(idResponse{ID: id})
}

type signInInput struct {
	PhoneNumber string `json:"phone_number" validate:"required"`
	Password    string `json:"password" validate:"required"`
}

type tokenResponse struct {
	AccessToken  string `json:"access"`
	RefreshToken string `json:"refresh"`
}

// @Tags auth
// @Description user sign in
// @ModuleID MemberSignIn
// @Accept  json
// @Produce  json
// @Param input body signInInput true "sign in info"
// @Success 200 {object} tokenResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /auth/sign-in [post]
func (h *Handler) userSignIn(c *fiber.Ctx) error {
	var input signInInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	ok, errs := validationStructs.ValidateStruct(input)

	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(errs)
	}
	user := domain.User{
		PhoneNumber: input.PhoneNumber,
		Password:    input.Password,
	}
	res, err := h.services.UserAuth.UserSignIn(user)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotExist) || errors.Is(err, domain.ErrUserBlocked) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(tokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
}

type TokenAuth struct {
	Token string `json:"token"`
}

type refreshInput struct {
	Token string `json:"token" binding:"required"`
}

// @Tags auth
// @Description refresh token
// @ModuleID refreshToken
// @Param input body refreshInput true "refresh info"
// @Success 200 {object} tokenResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /auth/refresh [post]
func (h *Handler) userRefresh(c *fiber.Ctx) error {
	var input refreshInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	res, err := h.services.UserAuth.RefreshToken(input.Token)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(tokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
}

// @Security User_Auth
// @Tags auth
// @Description get users info
// @ModuleId userMe
// @Accept  json
// @Produce  json
// @Success 200 {object} domain.UserList
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /auth/user/me [get]
func (h *Handler) userMe(c *fiber.Ctx) error {
	_, id := getUser(c)

	list, err := h.services.UserAuth.UserMe(id)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(list)
}

// @Tags auth
// @Description reset password
// @ModuleID resetPassword
// @Accept json
// @Produce json
// @Param data body domain.EmailInput true "reset password via Email"
// @Success 201 {object} TokenAuth
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /auth/verify [post]
func (h *Handler) verifyEmail(c *fiber.Ctx) error {
	var input domain.EmailInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	ok, errs := validationStructs.ValidateStruct(input)

	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(errs)
	}

	host := parseRequestHost(c)

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	token, err := h.services.UserAuth.VerifyEmail(input.Email, host)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotExist) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(TokenAuth{Token: token})

}

type ResetPasswordInput struct {
	Token              string `json:"token" validate:"required"`
	NewPassword        string `json:"new_password"  validate:"required,min=8,max=64"`
	NewPasswordConfirm string `json:"new_password_confirm" validate:"required"`
}

// @Tags auth
// @Description reset password
// @ModuleID resetPassword
// @Accept json
// @Produce json
// @Param data body ResetPasswordInput true "reset password "
// @Success 201 {object} okResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /auth/reset-password [post]
func (h *Handler) resetPassword(c *fiber.Ctx) error {
	var input ResetPasswordInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	id, err := h.userIdentity(input.Token)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	err = h.services.UserAuth.ResetPassword(id, input.NewPassword, input.NewPasswordConfirm)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(okResponse{Message: "OK"})

}

// @Tags auth
// @Description gets all users
// @ID getAllUsers
// @Accept json
// @Produce json
// @Param array query domain.Pagination  true "A page info"
// @Success 200 {object} domain.GetAllUsersCategoryResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /auth/users [get]
func (h *Handler) getAllUsers(c *fiber.Ctx) error {
	var page domain.Pagination
	if err := c.QueryParser(&page); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	list, err := h.services.UserAuth.GetAll(page)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(list)
}

// @Tags auth
// @Description download members to excel
// @ID downloadUsers
// @Produce application/octet-stream
// @Success 200 {file} file binary
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /auth/download/users [get]
func (h *Handler) downloadUsers(c *fiber.Ctx) error {
	cellValue := map[string]string{"A1": "id", "B1": "имя", "C1": "номер телефона", "D1": "email", "E1": "пол", "F1": "возраст", "G1": "город", "H1": "дата записи", "I1": "статус"}

	file, err := excel.File(cellValue)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}

	file, err = h.services.UserAuth.DownloadUsers(file)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}

	var b bytes.Buffer

	if err := file.Write(&b); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	downloadName := time.Now().UTC().Format("users.xlsx")
	c.GetRespHeader("Content-Description", "File Transfer")
	c.GetRespHeader("Content-Type", "application/octet-stream")
	c.GetRespHeader("Content-Disposition", "attachment; filename="+downloadName)

	err = file.SaveAs(downloadName)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}

	return c.Download(downloadName)

}

// @Security User_Auth
// @Tags auth
// @Description block users
// @ModuleID blockUser
// @Accept json
// @Produce json
// @Param id path string true "user id"
// @Success 200 {object} okResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /auth/block-user/{id} [put]
func (h *Handler) blockUser(c *fiber.Ctx) error {

	userType, _ := getUser(c)

	if userType != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(response{Message: "нет доступа"})
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	if err = h.services.UserAuth.BlockUser(id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(okResponse{Message: "OK"})
}

// @Security User_Auth
// @Tags auth
// @Description unblock users
// @ModuleID unblockUser
// @Accept json
// @Produce json
// @Param id path string true "user id"
// @Success 200 {object} okResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /auth/unblock-user/{id} [put]
func (h *Handler) unblockUser(c *fiber.Ctx) error {

	userType, _ := getUser(c)

	if userType != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(response{Message: "нет доступа"})
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	if err = h.services.UserAuth.UnBlockUser(id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(okResponse{Message: "OK"})

}
