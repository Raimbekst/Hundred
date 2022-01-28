package v1

import (
	"HundredToFive/internal/domain"
	"errors"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"strconv"
)

func (h *Handler) initConditionRoutes(api fiber.Router) {
	desc := api.Group("/condition")
	{

		desc.Get("/", h.getAllConditions)
		desc.Get("/:id", h.getConditionById)
		admin := desc.Group("", jwtware.New(jwtware.Config{
			SigningKey: []byte(h.signingKey),
		}))
		{
			admin.Post("", h.createCondition)
			admin.Post("/kz", h.createConditionKz)
			admin.Put("/:id", h.updateCondition)
			admin.Delete("/:id", h.deleteCondition)
		}
	}
}

type Condition struct {
	Caption string `json:"caption" db:"caption"`
	Text    string `json:"text" db:"text"`
}

// @Security User_Auth
// @Tags condition
// @ModuleID createConditionKz
// @Accept  json
// @Produce  json
// @Param data body Condition  true "condition"
// @Success 201 {object} idResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /condition/kz [post]
func (h *Handler) createConditionKz(c *fiber.Ctx) error {
	userType, _ := getUser(c)

	if userType != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(response{Message: "нет доступа"})
	}

	var input Condition
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	con := domain.Condition{
		Caption:      input.Caption,
		Text:         input.Text,
		LanguageType: "kz",
	}

	id, err := h.services.Condition.Create(con)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(idResponse{ID: id})
}

// @Security User_Auth
// @Tags condition
// @ModuleID createCondition
// @Accept  json
// @Produce  json
// @Param data body Condition  true "condition"
// @Success 201 {object} idResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /condition [post]
func (h *Handler) createCondition(c *fiber.Ctx) error {
	userType, _ := getUser(c)

	if userType != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(response{Message: "нет доступа"})
	}

	var input Condition
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	con := domain.Condition{
		Caption:      input.Caption,
		Text:         input.Text,
		LanguageType: "ru",
	}

	id, err := h.services.Condition.Create(con)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(idResponse{ID: id})
}

// @Tags condition
// @Description get all conditions
// @ID get-all-conditions
// @Accept  json
// @Produce  json
// @Param array query domain.Pagination  true "A page info"
// @Param filter query LanguageTypeInput true "A filter info"
// @Success 200 {object} domain.GetAllDescCategoryResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /condition [get]
func (h *Handler) getAllConditions(c *fiber.Ctx) error {
	var lang LanguageTypeInput
	if err := c.QueryParser(&lang); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	var page domain.Pagination
	if err := c.QueryParser(&page); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	list, err := h.services.Condition.GetAll(page, lang.LanguageType)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(list)

}

// @Tags condition
// @Description get condition by id
// @ID get-condition-by-id
// @Accept  json
// @Produce  json
// @Param id path string true "condition id"
// @Success 200 {object} domain.Description
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /condition/{id} [get]
func (h *Handler) getConditionById(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	list, err := h.services.Condition.GetById(id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(list)
}

// @Security User_Auth
// @Tags condition
// @Description  update  condition
// @ModuleID updateCondition
// @Accept  json
// @Produce  json
// @Param id path string true "condition id"
// @Param input body Condition false "condition"
// @Success 200 {object} okResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /condition/{id} [put]
func (h *Handler) updateCondition(c *fiber.Ctx) error {
	userType, _ := getUser(c)

	if userType != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(response{Message: "нет доступа"})
	}
	var input Condition
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	con := domain.Condition{
		Caption: input.Caption,
		Text:    input.Text,
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	if err := h.services.Condition.Update(id, con); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(okResponse{Message: "OK"})

}

// @Security User_Auth
// @Tags condition
// @Description delete condition
// @ModuleID deleteCondition
// @Accept  json
// @Produce  json
// @Param id path string true "condition id"
// @Success 200 {object} okResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /condition/{id} [delete]
func (h *Handler) deleteCondition(c *fiber.Ctx) error {
	userType, _ := getUser(c)

	if userType != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(response{Message: "нет доступа"})
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	if err := h.services.Condition.Delete(id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(okResponse{Message: "OK"})
}
