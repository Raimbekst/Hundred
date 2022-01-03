package v1

import (
	"HundredToFive/internal/domain"
	"errors"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"strconv"
)

func (h *Handler) initCityCategoryRoutes(api fiber.Router) {
	city := api.Group("/city")
	{

		city.Get("/", h.getAllCities)
		city.Get("/:id", h.getCityById)
		admin := city.Group("", jwtware.New(jwtware.Config{
			SigningKey: []byte(h.signingKey),
		}))
		{
			admin.Post("", h.createCity)
			admin.Put("/:id", h.updateCity)
			admin.Delete("/:id", h.deleteCity)
		}
	}
}

type City struct {
	Name string `json:"name" validate:"required"`
}

// @Security User_Auth
// @Tags city
// @ModuleID createCity
// @Accept  json
// @Produce  json
// @Param data body City  true "city"
// @Success 201 {object} idResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /city [post]
func (h *Handler) createCity(c *fiber.Ctx) error {
	userType, _ := getUser(c)

	if userType != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(response{Message: "нет доступа"})
	}

	var input City
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	city := domain.City{
		Name: input.Name,
	}
	id, err := h.services.City.Create(city)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(idResponse{ID: id})
}

// @Tags city
// @Description get all cities
// @ID get-all-cities
// @Accept  json
// @Produce  json
// @Param array query domain.Pagination  true "A page info"
// @Success 200 {object} domain.GetAllCityCategoryResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /city [get]
func (h *Handler) getAllCities(c *fiber.Ctx) error {
	var page domain.Pagination
	if err := c.QueryParser(&page); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	list, err := h.services.City.GetAll(page)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(list)

}

// @Tags city
// @Description get city by id
// @ID get-city-by-id
// @Accept  json
// @Produce  json
// @Param id path string true "city id"
// @Success 200 {object} domain.City
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /city/{id} [get]
func (h *Handler) getCityById(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	list, err := h.services.City.GetById(id)
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
func (h *Handler) updateCity(c *fiber.Ctx) error {
	userType, _ := getUser(c)

	if userType != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(response{Message: "нет доступа"})
	}
	var input City
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	city := domain.City{
		Name: input.Name,
	}
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	if err := h.services.City.Update(id, city); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(okResponse{Message: "OK"})

}

// @Security User_Auth
// @Tags city
// @Description delete city
// @ModuleID deleteCity
// @Accept  json
// @Produce  json
// @Param id path string true "city id"
// @Success 200 {object} okResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /city/{id} [delete]
func (h *Handler) deleteCity(c *fiber.Ctx) error {
	userType, _ := getUser(c)

	if userType != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(response{Message: "нет доступа"})
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	if err := h.services.City.Delete(id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(okResponse{Message: "OK"})
}
