package v1

import (
	"HundredToFive/internal/domain"
	"HundredToFive/pkg/validation/validationStructs"
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
	raffle := api.Group("/notification")
	{
		raffle.Get("/", h.getAllNotifications)
		raffle.Get("/:id", h.getNotificationById)
		raffle.Post("/response", h.notificationResponse)
		admin := raffle.Group("", jwtware.New(
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
	Time      float32 `json:"time"`
}

// @Security User_Auth
// @Tags notification
// @Description create notification for all users
// @ID createNotyForAllUsers
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

	if err := c.QueryParser(&input); err != nil {
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
		Status:    planned,
		Getters:   getterAll,
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

	// This registration token comes from the client FCM SDKs.
	registrationToken := "YOUR_REGISTRATION_TOKEN"

	// See documentation on defining a message payload.
	message := &messaging.Message{
		Data: map[string]string{
			"score": "850",
			"time":  "2:45",
		},
		Token: registrationToken,
	}

	// Send a message to the device corresponding to the provided
	// registration token.
	res, err := cl.Send(c.Context(), message)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	// Response is a message ID string.
	fmt.Println("Successfully sent message:", res)

	fmt.Println(app)

	var input NotificationForUser

	if err := c.QueryParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	ok, errs := validationStructs.ValidateStruct(input)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(errs)
	}

	noty := domain.Notification{
		Title:   input.Title,
		Text:    input.Text,
		Link:    input.Link,
		Status:  sent,
		Getters: getterUser,
		Date:    float64(time.Now().Unix()),
		Ids:     input.UserIds,
	}

	_, id, err := h.services.Notification.CreateForUser(noty)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(id)

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
	var page domain.Pagination
	if err := c.QueryParser(&page); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	list, err := h.services.Notification.GetAll(page)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(list)

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
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	list, err := h.services.Notification.GetById(id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(list)
}

// @Security User_Auth
// @Tags city
// @Description  update  city
// @ModuleID updateCity
// @Accept  json
// @Produce  json
// @Param id path string true "city id"
// @Param input body City false "city"
// @Success 200 {object} okResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /city/{id} [put]
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

	var input RegistrationTokenInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	ok, errs := validationStructs.ValidateStruct(input)

	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(errs)
	}

	id, err := h.services.Notification.StoreUsersToken(input.RegistrationToken)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(idResponse{ID: id})
}
