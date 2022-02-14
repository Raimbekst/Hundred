package v1

import (
	"HundredToFive/internal/domain"
	"HundredToFive/pkg/excel"
	"HundredToFive/pkg/validation/validationStructs"
	"bytes"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"strconv"
	"strings"
	"time"
)

func (h *Handler) initNotificationRoutes(api fiber.Router) {
	noty := api.Group("/notification")
	{
		noty.Get("/", h.getAllNotifications)
		noty.Get("/download", h.downloadNotification)
		noty.Get("/:id", h.getNotificationById)
		noty.Post("/response", h.notificationResponse)
		admin := noty.Group("", jwtware.New(
			jwtware.Config{
				SigningKey: []byte(h.signingKey),
			}))
		{
			admin.Post("/user", h.createNotyForSpecificUser)
			admin.Post("/all", h.createNotyForAllUsers)
			admin.Put("/:id", h.updateNotification)
			admin.Delete("/:id", h.deleteNotification)
		}
	}
}

type Notification struct {
	Title     string  `json:"title" validate:"required"`
	Text      string  `json:"text"  validate:"required"`
	PartnerId int     `json:"partner_id"`
	Link      string  `json:"link" `
	Reference string  `json:"reference" `
	Date      float64 `json:"date"`
	Time      int     `json:"time"`
}

// @Security User_Auth
// @Tags notification
// @Description create notification for all users
// @ModuleID createNotyForAllUsers
// @Accept json
// @Produce json
// @Param data body Notification true "notification"
// @Success 201 {object} idResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /notification/all [post]
func (h *Handler) createNotyForAllUsers(c *fiber.Ctx) error {
	url := c.BaseURL()

	userType, _ := getUser(c)

	if userType != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(response{Message: "нет доступа"})
	}

	var input Notification

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	ok, errs := validationStructs.ValidateStruct(input)

	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(errs)
	}

	partnerList, _ := h.services.Partner.GetById(input.PartnerId)

	var logo string
	if partnerList.Logo != "" {
		logo = url + "/" + "media/" + partnerList.Logo
	}

	var partId *int = &input.PartnerId

	if input.PartnerId == 0 {
		partId = nil
	}

	notyDate := input.Date + float64(input.Time)

	noty := domain.Notification{
		Title:     input.Title,
		Text:      input.Text,
		PartnerId: partId,
		Link:      input.Link,
		Reference: input.Reference,
		Date:      notyDate,
		Time:      input.Time,
		Status:    1,
		Getters:   1,
		Logo:      logo,
	}
	id, err := h.services.Notification.Create(noty)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(idResponse{ID: id})
}

type NotificationForUser struct {
	UserIds []int  `json:"user_ids" validate:"required"`
	Title   string `json:"title"  validate:"required"`
	Text    string `json:"text" validate:"required"`
	Link    string `json:"link" `
}

type MessageSent struct {
	Id       int `json:"id"`
	Response int `json:"response"`
}

// @Security User_Auth
// @Tags notification
// @Description create notification for specific user
// @ID createNotyForSpecificUser
// @Accept json
// @Produce json
// @Param data body NotificationForUser true "notification"
// @Success 201 {object} idResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /notification/user [post]
func (h *Handler) createNotyForSpecificUser(c *fiber.Ctx) error {
	userType, _ := getUser(c)

	if userType != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(response{Message: "нет доступа"})
	}

	var input NotificationForUser

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	ok, errs := validationStructs.ValidateStruct(input)

	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(errs)
	}

	hour := (time.Now().Hour() + 6) * 3600

	minute := time.Now().Minute() * 60

	inp := domain.Notification{
		Title:   input.Title,
		Text:    input.Text,
		Link:    input.Link,
		Date:    float64(time.Now().Unix()),
		Time:    hour + minute,
		Status:  1,
		Getters: 2,
		Ids:     input.UserIds,
	}

	tokens, idInt, err := h.services.Notification.CreateForUser(inp)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}

	if tokens == nil {
		tokens = []string{"Random"}
	}
	res, err := h.firebaseNotification(h.ctx, inp, tokens, idInt)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(MessageSent{Id: idInt, Response: res.SuccessCount})

}

// @Tags notification
// @Description get all notifications
// @ID get-all-notifications
// @Accept  json
// @Produce  json
// @Param array query domain.Pagination  true "A page info"
// @Success 200 {object} domain.GetAllNotificationsResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /notification [get]
func (h *Handler) getAllNotifications(c *fiber.Ctx) error {
	url := c.BaseURL()

	var page domain.Pagination

	if err := c.QueryParser(&page); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	list, err := h.services.Notification.GetAll(page)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}

	for _, value := range list.Data {
		if value.Logo != "" {
			value.Logo = url + "/" + "media/" + value.Logo
		}
	}
	return c.Status(fiber.StatusOK).JSON(list)

}

// @Tags notification
// @Description download notifications to excel
// @ID downloadNotification
// @Produce application/octet-stream
// @Success 200 {file} file binary
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /notification/download [get]
func (h *Handler) downloadNotification(c *fiber.Ctx) error {

	cellValue := map[string]string{"A1": "id", "B1": "title", "C1": "текст", "D1": "статус", "E1": "ссылка", "F1": "примечение", "G1": "получатели"}

	file, err := excel.File(cellValue)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}

	file, err = h.services.Notification.DownloadNotification(file)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}

	var b bytes.Buffer

	if err := file.Write(&b); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	downloadName := time.Now().UTC().Format("notifications.xlsx")
	c.GetRespHeader("Content-Description", "File Transfer")
	c.GetRespHeader("Content-Type", "application/octet-stream")
	c.GetRespHeader("Content-Disposition", "attachment; filename="+downloadName)

	err = file.SaveAs(downloadName)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}

	return c.Download(downloadName)
}

// @Tags notification
// @Description get notification by id
// @ID get-notification-by-id
// @Accept  json
// @Produce  json
// @Param id path string true "notification id"
// @Success 200 {object} domain.Notification
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /notification/{id} [get]
func (h *Handler) getNotificationById(c *fiber.Ctx) error {
	url := c.BaseURL()
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	fmt.Println(id)
	list, err := h.services.Notification.GetById(id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}

	if list.Logo != "" {
		list.Logo = url + "/" + "media/" + list.Logo
	}

	return c.Status(fiber.StatusOK).JSON(list)
}

// @Security User_Auth
// @Tags notification
// @Description  update  notification
// @ModuleID updateNotification
// @Accept  json
// @Produce  json
// @Param id path string true "notification id"
// @Param input body Notification false "notification"
// @Success 200 {object} okResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /notification/{id} [put]
func (h *Handler) updateNotification(c *fiber.Ctx) error {
	userType, _ := getUser(c)

	if userType != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(response{Message: "нет доступа"})
	}
	var input Notification

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	notyDate := input.Date + float64(input.Time)

	noty := domain.Notification{
		Title:     input.Title,
		Text:      input.Text,
		PartnerId: &input.PartnerId,
		Link:      input.Link,
		Reference: input.Reference,
		Date:      notyDate,
		Time:      input.Time,
	}

	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	fmt.Println(id)
	if err := h.services.Notification.Update(id, noty); err != nil {
		if errors.Is(err, domain.ErrNotFound) || errors.Is(err, domain.ErrUpdateNotification) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(okResponse{Message: "OK"})

}

// @Security User_Auth
// @Tags notification
// @Description delete notification
// @ModuleID deleteNotification
// @Accept  json
// @Produce  json
// @Param id path string true "notification id"
// @Success 200 {object} okResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /notification/{id} [delete]
func (h *Handler) deleteNotification(c *fiber.Ctx) error {

	userType, _ := getUser(c)

	if userType != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(response{Message: "нет доступа"})
	}
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	if err := h.services.Notification.Delete(id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(okResponse{Message: "OK"})

}

type RegistrationTokenInput struct {
	RegistrationToken string `json:"registration_token" validate:"required"`
}

// @Tags notification
// @Description accept registration token for push notification
// @ID notification-responsible
// @Accept json
// @Produce json
// @Param input body RegistrationTokenInput true "registration token"
// @Success 200 {object} okResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /notification/response [post]
func (h *Handler) notificationResponse(c *fiber.Ctx) error {

	var userId *int = nil
	header := string(c.Request().Header.Peek("Authorization"))

	if header != "" {
		headerParts := strings.Split(header, " ")
		uId, _, err := h.tokenManager.Parse(headerParts[1])
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: "token can not get user info"})
		}

		idInt, err := strconv.Atoi(uId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: "can not convert string to int"})
		}
		userId = &idInt
	}

	var input RegistrationTokenInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	ok, errs := validationStructs.ValidateStruct(input)

	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(errs)
	}

	id, err := h.services.Notification.StoreUsersToken(userId, input.RegistrationToken)

	if err != nil {
		if errors.Is(err, domain.ErrTokenAlreadyExist) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(idResponse{ID: id})
}
