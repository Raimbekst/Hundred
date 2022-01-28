package v1

import (
	"HundredToFive/internal/domain"
	"errors"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"strconv"
)

type LanguageTypeInput struct {
	LanguageType string `json:"language_type" query:"language_type" form:"language_type" enums:"ru,kz"`
}

func (h *Handler) initFaqCategoryRoutes(api fiber.Router) {
	city := api.Group("/faq")
	{
		city.Get("/", h.getAllFaqs)
		city.Get("/:id", h.getFaqById)
		admin := city.Group("", jwtware.New(jwtware.Config{
			SigningKey: []byte(h.signingKey),
		}))
		{
			admin.Post("", h.createFaq)
			admin.Post("/kz", h.createFaqKz)
			admin.Put("/:id", h.updateFaq)
			admin.Delete("/:id", h.deleteFaq)
		}
	}
}

func (h *Handler) initDescCategoryRoutes(api fiber.Router) {
	desc := api.Group("/description")
	{

		desc.Get("/", h.getAllDescriptions)
		desc.Get("/:id", h.getDescById)
		admin := desc.Group("", jwtware.New(jwtware.Config{
			SigningKey: []byte(h.signingKey),
		}))
		{
			admin.Post("", h.createDesc)
			admin.Post("/kz", h.createDescKz)
			admin.Put("/:id", h.updateDesc)
			admin.Delete("/:id", h.deleteDesc)
		}
	}
}

type Faq struct {
	Question string `json:"question" db:"question"`
	Answer   string `json:"answer" db:"answer"`
}

// @Security User_Auth
// @Tags faq
// @ModuleID createFaqKz
// @Accept  json
// @Produce  json
// @Param data body Faq  true "faq"
// @Success 201 {object} idResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /faq/kz [post]
func (h *Handler) createFaqKz(c *fiber.Ctx) error {
	userType, _ := getUser(c)

	if userType != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(response{Message: "нет доступа"})
	}

	var input Faq
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	faq := domain.Faq{
		Answer:       input.Answer,
		Question:     input.Question,
		LanguageType: "kz",
	}
	id, err := h.services.Faq.Create(faq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(idResponse{ID: id})
}

// @Security User_Auth
// @Tags faq
// @ModuleID createFaq
// @Accept  json
// @Produce  json
// @Param data body Faq  true "faq"
// @Success 201 {object} idResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /faq [post]
func (h *Handler) createFaq(c *fiber.Ctx) error {
	userType, _ := getUser(c)

	if userType != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(response{Message: "нет доступа"})
	}

	var input Faq
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	faq := domain.Faq{
		Answer:       input.Answer,
		Question:     input.Question,
		LanguageType: "ru",
	}
	id, err := h.services.Faq.Create(faq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(idResponse{ID: id})
}

// @Tags faq
// @Description get all faqs
// @ID get-all-faqs
// @Accept  json
// @Produce  json
// @Param array query domain.Pagination  true "A page info"
// @Param filter query LanguageTypeInput true "A filter info"
// @Success 200 {object} domain.GetAllFaqsCategoryResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /faq [get]
func (h *Handler) getAllFaqs(c *fiber.Ctx) error {
	var page domain.Pagination
	if err := c.QueryParser(&page); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	var filter LanguageTypeInput

	if err := c.QueryParser(&filter); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	list, err := h.services.Faq.GetAll(page, filter.LanguageType)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(list)

}

// @Tags faq
// @Description get faq by id
// @ID get-faq-by-id
// @Accept  json
// @Produce  json
// @Param id path string true "faq id"
// @Success 200 {object} domain.Faq
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /faq/{id} [get]
func (h *Handler) getFaqById(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	list, err := h.services.Faq.GetById(id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(list)
}

// @Security User_Auth
// @Tags faq
// @Description  update  faq
// @ModuleID updateFaq
// @Accept  json
// @Produce  json
// @Param id path string true "faq id"
// @Param input body Faq false "faq"
// @Success 200 {object} okResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /faq/{id} [put]
func (h *Handler) updateFaq(c *fiber.Ctx) error {
	userType, _ := getUser(c)

	if userType != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(response{Message: "нет доступа"})
	}
	var input Faq
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	faq := domain.Faq{
		Answer:   input.Answer,
		Question: input.Question,
	}
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	if err := h.services.Faq.Update(id, faq); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(okResponse{Message: "OK"})

}

// @Security User_Auth
// @Tags faq
// @Description delete faq
// @ModuleID deleteCity
// @Accept  json
// @Produce  json
// @Param id path string true "faq id"
// @Success 200 {object} okResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /faq/{id} [delete]
func (h *Handler) deleteFaq(c *fiber.Ctx) error {
	userType, _ := getUser(c)

	if userType != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(response{Message: "нет доступа"})
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	if err := h.services.Faq.Delete(id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(okResponse{Message: "OK"})
}

////////////////////////////////////////////////////

type Description struct {
	Caption string `json:"caption" db:"caption"`
	Text    string `json:"text" db:"text"`
}

// @Security User_Auth
// @Tags about
// @ModuleID createDescKz
// @Accept  json
// @Produce  json
// @Param data body Description  true "description"
// @Success 201 {object} idResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /description/kz [post]
func (h *Handler) createDescKz(c *fiber.Ctx) error {
	userType, _ := getUser(c)

	if userType != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(response{Message: "нет доступа"})
	}

	var input Description
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	desc := domain.Description{
		Caption:      input.Caption,
		Text:         input.Text,
		LanguageType: "kz",
	}

	id, err := h.services.Faq.CreateDesc(desc)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(idResponse{ID: id})
}

// @Security User_Auth
// @Tags about
// @ModuleID createDesc
// @Accept  json
// @Produce  json
// @Param data body Description  true "description"
// @Success 201 {object} idResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /description [post]
func (h *Handler) createDesc(c *fiber.Ctx) error {
	userType, _ := getUser(c)

	if userType != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(response{Message: "нет доступа"})
	}

	var input Description
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	desc := domain.Description{
		Caption:      input.Caption,
		Text:         input.Text,
		LanguageType: "ru",
	}

	id, err := h.services.Faq.CreateDesc(desc)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(idResponse{ID: id})
}

// @Tags about
// @Description get all description
// @ID get-all-description
// @Accept  json
// @Produce  json
// @Param array query domain.Pagination  true "A page info"
// @Param filter query LanguageTypeInput true "A filter info"
// @Success 200 {object} domain.GetAllDescCategoryResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /description [get]
func (h *Handler) getAllDescriptions(c *fiber.Ctx) error {

	var lang LanguageTypeInput
	if err := c.QueryParser(&lang); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	var page domain.Pagination
	if err := c.QueryParser(&page); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	list, err := h.services.Faq.GetAllDesc(page, lang.LanguageType)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(list)

}

// @Tags about
// @Description get description by id
// @ID get-description-by-id
// @Accept  json
// @Produce  json
// @Param id path string true "description id"
// @Success 200 {object} domain.Description
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /description/{id} [get]
func (h *Handler) getDescById(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	list, err := h.services.Faq.GetDescById(id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(list)
}

// @Security User_Auth
// @Tags about
// @Description  update  description
// @ModuleID updateDesc
// @Accept  json
// @Produce  json
// @Param id path string true "description id"
// @Param input body Description false "description"
// @Success 200 {object} okResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /description/{id} [put]
func (h *Handler) updateDesc(c *fiber.Ctx) error {
	userType, _ := getUser(c)

	if userType != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(response{Message: "нет доступа"})
	}
	var input Description
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	desc := domain.Description{
		Caption: input.Caption,
		Text:    input.Text,
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	if err := h.services.Faq.UpdateDesc(id, desc); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(okResponse{Message: "OK"})

}

// @Security User_Auth
// @Tags about
// @Description delete description
// @ModuleID deleteDesc
// @Accept  json
// @Produce  json
// @Param id path string true "description id"
// @Success 200 {object} okResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /description/{id} [delete]
func (h *Handler) deleteDesc(c *fiber.Ctx) error {
	userType, _ := getUser(c)

	if userType != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(response{Message: "нет доступа"})
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	if err := h.services.Faq.DeleteDesc(id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(okResponse{Message: "OK"})
}
