package v1

import (
	"HundredToFive/internal/domain"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func (h *Handler) initAboutWebsiteRoutes(api fiber.Router) {
	banner := api.Group("/website")
	{
		banner.Get("", h.getAllInfo)
		admin := banner.Group("", jwtware.New(jwtware.Config{
			SigningKey: []byte(h.signingKey),
		}))
		{
			admin.Post("", h.createInfo)
			admin.Put("", h.updateInfo)
			admin.Delete("", h.deleteInfo)

		}
	}
}

// @Security User_Auth
// @Tags about
// @ModuleID createInfo
// @Accept  json
// @Produce  json
// @Param data body domain.AboutUs  true "about website"
// @Success 201 {object} idResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /website [post]
func (h *Handler) createInfo(c *fiber.Ctx) error {
	userType, _ := getUser(c)

	if userType != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(response{Message: "нет доступа"})
	}

	var input domain.AboutUs

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	err := h.services.About.Create(input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(okResponse{Message: "OK,Created"})
}

// @Tags about
// @Description get all website info
// @ID get-all-website-info
// @Accept  json
// @Produce  json
// @Success 200 {object} domain.AboutUs
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /website [get]
func (h *Handler) getAllInfo(c *fiber.Ctx) error {
	var page domain.Pagination
	if err := c.QueryParser(&page); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	list, err := h.services.About.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(list)

}

// @Security User_Auth
// @Tags about
// @Description  update  info
// @ModuleID updateInfo
// @Accept  json
// @Produce  json
// @Param input body domain.AboutUs false "website info"
// @Success 200 {object} okResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /website [put]
func (h *Handler) updateInfo(c *fiber.Ctx) error {
	userType, _ := getUser(c)

	if userType != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(response{Message: "нет доступа"})
	}

	var input domain.AboutUs
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	if err := h.services.About.Update(input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(okResponse{Message: "OK"})

}

// @Security User_Auth
// @Tags about
// @Description delete about
// @ModuleID deleteInfo
// @Accept  json
// @Produce  json
// @Success 200 {object} okResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /website [delete]
func (h *Handler) deleteInfo(c *fiber.Ctx) error {
	userType, _ := getUser(c)

	if userType != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(response{Message: "нет доступа"})
	}

	if err := h.services.About.Delete(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(okResponse{Message: "OK"})
}
