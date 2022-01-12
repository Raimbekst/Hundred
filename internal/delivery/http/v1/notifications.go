package v1

import (
	"HundredToFive/internal/domain"
	"HundredToFive/pkg/excel"
	"HundredToFive/pkg/validation/validationStructs"
	"bytes"
	"errors"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"fmt"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"google.golang.org/api/option"
	"strconv"
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

	noty := domain.Notification{
		Title:     input.Title,
		Text:      input.Text,
		PartnerId: input.PartnerId,
		Link:      input.Link,
		Reference: input.Reference,
		Date:      input.Date,
		Time:      input.Time,
		Status:    1,
		Getters:   1,
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
	Id       int                      `json:"id"`
	Response *messaging.BatchResponse `json:"response"`
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

	f := c.Get("Authorization")
	fmt.Println(f)
	var input NotificationForUser

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	ok, errs := validationStructs.ValidateStruct(input)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(errs)
	}

	opt := option.WithCredentialsFile("./hundredtofive-1652d-firebase-adminsdk-oecbb-6f47f9aa54.json")

	config := &firebase.Config{ProjectID: "hundredtofive-1652d"}

	app, err := firebase.NewApp(c.Context(), config, opt)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	cl, err := app.Messaging(c.Context())

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	registrationToken := []string{"e6PDZz6B-fficwwtO6ePWy:APA91bHpqyI084W-jio7_D_wRT9wi3WsWL2bg4p0UVCt2KumkAIuHRmlX3Wc6CFYzsaWoUROps3Y5PNtvbQjIbw_hHTuNmVKQ5A76_3s-IyEOrQsRJcIMhqD5UQngpkNgAz6FeygJ0-4"}

	message := &messaging.MulticastMessage{
		Data: map[string]string{
			"title": input.Title,
			"text":  input.Text,
			"link":  input.Link,
		},
		Tokens: registrationToken,
	}

	res, err := cl.SendMulticast(c.Context(), message)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(MessageSent{Response: res})

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
	noty := domain.Notification{
		Title:     input.Title,
		Text:      input.Text,
		PartnerId: input.PartnerId,
		Link:      input.Link,
		Reference: input.Reference,
		Date:      input.Date,
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

	var id int = 0

	str := c.Get("")

	fmt.Println(str)

	var input RegistrationTokenInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	ok, errs := validationStructs.ValidateStruct(input)

	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(errs)
	}

	id, err := h.services.Notification.StoreUsersToken(&id, input.RegistrationToken)

	if err != nil {
		if errors.Is(err, domain.ErrTokenAlreadyExist) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(idResponse{ID: id})
}
