package v1

import (
	"HundredToFive/internal/domain"
	"HundredToFive/pkg/media"
	"errors"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"mime/multipart"
	"strconv"
)

func (h *Handler) initBannerCategoryRoutes(api fiber.Router) {
	banner := api.Group("/banner")
	{
		banner.Get("/", h.getAllBanners)
		banner.Get("/:id", h.getBannerById)

		admin := banner.Group("", jwtware.New(
			jwtware.Config{
				SigningKey: []byte(h.signingKey),
			}))

		{
			admin.Post("", h.createBanner)
			admin.Post("/kz", h.createBannerKz)
			admin.Put("/:id", h.updateBanner)
			admin.Delete("/:id", h.deleteBanner)

		}

	}
}

type Banner struct {
	Name   string                `form:"name"`
	Status int                   `form:"status"  enums:"1,2" default:"1"`
	Image  *multipart.FileHeader `form:"image"`
	Iframe string                `form:"iframe"`
}

// @Security User_Auth
// @Tags banner
// @ModuleID createBannerKz
// @Accept  multipart/form-data
// @Produce  json
// @Param name formData string false "name of banner"
// @Param status formData int false "only 1 or 2" Enums(1,2)
// @Param image formData file false "image"
// @Param iframe formData string false "iframe"
// @Success 201 {object} idResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /banner/kz [post]
func (h *Handler) createBannerKz(c *fiber.Ctx) error {

	userType, _ := getUser(c)

	if userType != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(response{Message: "нет доступа"})
	}
	var (
		input Banner
		err   error
	)

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	var image string

	file, _ := c.FormFile("image")

	if file != nil {
		image, err = media.GetFileName(c, file)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
	}

	banner := domain.Banner{
		Name:         input.Name,
		Status:       input.Status,
		Image:        image,
		Iframe:       input.Iframe,
		LanguageType: "kz",
	}
	id, err := h.services.Banner.Create(banner)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(idResponse{ID: id})
}

// @Security User_Auth
// @Tags banner
// @ModuleID createBanner
// @Accept  multipart/form-data
// @Produce  json
// @Param name formData string false "name of banner"
// @Param status formData int false "only 1 or 2" Enums(1,2)
// @Param image formData file false "image"
// @Param iframe formData string false "iframe"
// @Success 201 {object} idResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /banner [post]
func (h *Handler) createBanner(c *fiber.Ctx) error {

	userType, _ := getUser(c)

	if userType != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(response{Message: "нет доступа"})
	}
	var (
		input Banner
		err   error
	)

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	var image string

	file, _ := c.FormFile("image")

	if file != nil {
		image, err = media.GetFileName(c, file)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
	}

	banner := domain.Banner{
		Name:         input.Name,
		Status:       input.Status,
		Image:        image,
		Iframe:       input.Iframe,
		LanguageType: "ru",
	}
	id, err := h.services.Banner.Create(banner)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(idResponse{ID: id})
}

type FilterForBanner struct {
	Status int `json:"status" form:"status" query:"status" enums:"1,2"`
}

// @Tags banner
// @Description get all banners
// @ID get-all-banners
// @Accept  json
// @Produce  json
// @Param filter query FilterForBanner false "filter by status"
// @Param filter query LanguageTypeInput true "A filter info"
// @Param array query domain.Pagination  true "A page info"
// @Success 200 {object} domain.GetAllBannersCategoryResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /banner [get]
func (h *Handler) getAllBanners(c *fiber.Ctx) error {
	url := c.BaseURL()
	var filter FilterForBanner

	if err := c.QueryParser(&filter); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	var lang LanguageTypeInput

	if err := c.QueryParser(&lang); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	var page domain.Pagination
	if err := c.QueryParser(&page); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	list, err := h.services.Banner.GetAll(page, filter.Status, lang.LanguageType)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}

	for _, value := range list.Data {
		if value.Image != "" {
			value.Image = url + "/" + "media/" + value.Image
		}
	}

	return c.Status(fiber.StatusOK).JSON(list)

}

// @Tags banner
// @Description get banner by id
// @ID get-banner-by-id
// @Accept  json
// @Produce  json
// @Param id path string true "banner id"
// @Success 200 {object} domain.Banner
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /banner/{id} [get]
func (h *Handler) getBannerById(c *fiber.Ctx) error {
	url := c.BaseURL()
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	list, err := h.services.Banner.GetById(id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	if list.Image != "" {
		list.Image = url + "/" + "media/" + list.Image
	}
	return c.Status(fiber.StatusOK).JSON(list)
}

// @Security User_Auth
// @Tags banner
// @Description  update  banner
// @ModuleID updateBanner
// @Accept  multipart/form-data
// @Produce  json
// @Param id path string true "banner id"
// @Param name formData string false "name of banner"
// @Param status formData int false "only 1 or 2" Enums(1,2)
// @Param image formData file false "image"
// @Param iframe formData string false "iframe"
// @Success 200 {object} okResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /banner/{id} [put]
func (h *Handler) updateBanner(c *fiber.Ctx) error {
	userType, _ := getUser(c)

	if userType != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(response{Message: "нет доступа"})
	}

	var (
		input Banner
		err   error
	)
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	var image string

	file, _ := c.FormFile("image")

	if file != nil {
		image, err = media.GetFileName(c, file)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
	}

	banner := domain.Banner{
		Name:   input.Name,
		Status: input.Status,
		Image:  image,
		Iframe: input.Iframe,
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	if err := h.services.Banner.Update(id, banner); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(okResponse{Message: "OK"})

}

// @Security User_Auth
// @Tags banner
// @Description delete banner
// @ModuleID deleteBanner
// @Accept  json
// @Produce  json
// @Param id path string true "banner id"
// @Success 200 {object} okResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /banner/{id} [delete]
func (h *Handler) deleteBanner(c *fiber.Ctx) error {
	userType, _ := getUser(c)

	if userType != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(response{Message: "нет доступа"})
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	if err := h.services.Banner.Delete(id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(okResponse{Message: "OK"})
}
