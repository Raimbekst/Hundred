package v1

import (
	"HundredToFive/internal/domain"
	"HundredToFive/pkg/validation/validationStructs"
	"errors"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

type (
	winnersQuery struct {
		Date int64 `json:"date" query:"date" form:"date"`
	}
	raffleFilterQuery struct {
		RaffleId int `json:"raffle_id" query:"raffle_id" form:"raffle_id"`
	}
	monthFilter struct {
		Month int `json:"month" query:"month" form:"month"`
	}
)

func (h *Handler) initWinnerCategoryRoutes(api fiber.Router) {
	api.Get("/members", h.getAllMembers)
	api.Get("/days", h.getAllDays)
	api.Get("/months", h.getAllMonths)
	winner := api.Group("/winner")
	{
		winner.Get("", h.getAllWinnersOfToday)
		admin := winner.Group("", jwtware.New(
			jwtware.Config{
				SigningKey: []byte(h.signingKey),
			}))
		{
			admin.Post("", h.createWinner)
		}

	}
}

// @Security User_Auth
// @Tags winners
// @ModuleID createWinner
// @Accept  json
// @Produce  json
// @Param data body domain.WinnerInput  true "winners"
// @Success 200 {object} okResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /winner [post]
func (h *Handler) createWinner(c *fiber.Ctx) error {
	userType, _ := getUser(c)

	if userType != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(response{Message: "нет доступа"})
	}

	var input domain.WinnerInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	ok, errs := validationStructs.ValidateStruct(input)

	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(errs)
	}

	err := h.services.Winner.CreateWinner(input)

	if err != nil {
		if errors.Is(err, domain.ErrWinnerAlreadyExistInRaffle) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(okResponse{Message: "OK"})

}

// @Tags winners
// @Description gets all winners
// @ID get-all-winners
// @Accept json
// @Produce json
// @Param filter query winnersQuery true "day winner"
// @Param array query domain.Pagination  true "A page info"
// @Success 200 {object} domain.GetAllWinnersCategoryResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /winner [get]
func (h *Handler) getAllWinnersOfToday(c *fiber.Ctx) error {
	var page domain.Pagination
	if err := c.QueryParser(&page); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	var filterDate winnersQuery

	if err := c.QueryParser(&filterDate); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	list, err := h.services.Winner.GetAll(page, filterDate.Date)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(list)
}

// @Tags winners
// @Description gets all members
// @ID getAllMembers
// @Accept json
// @Produce json
// @Param filter query raffleFilterQuery true "get members"
// @Success 200 {object} domain.GetAllWinnersCategoryResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /members [get]
func (h *Handler) getAllMembers(c *fiber.Ctx) error {
	var page domain.Pagination
	if err := c.QueryParser(&page); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	var raffleFilter raffleFilterQuery

	if err := c.QueryParser(&raffleFilter); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	list, err := h.services.Winner.GetAllMembers(page, raffleFilter.RaffleId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(list)
}

// @Tags winners
// @Description gets all days
// @ID getAllDays
// @Accept json
// @Produce json
// @Param filter query monthFilter true "month for get winners"
// @Success 200 {object} domain.GetAllDaysResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /days [get]
func (h *Handler) getAllDays(c *fiber.Ctx) error {
	var page domain.Pagination
	if err := c.QueryParser(&page); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	var filterID monthFilter

	if err := c.QueryParser(&filterID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	list, err := h.services.Winner.GetAllDays(page, filterID.Month)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(list)
}

// @Tags winners
// @Description gets all days
// @ID getAllMonths
// @Accept json
// @Produce json
// @Param array query domain.Pagination  true "A page info"
// @Success 200 {object} domain.GetAllDaysResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /months [get]
func (h *Handler) getAllMonths(c *fiber.Ctx) error {
	var page domain.Pagination
	if err := c.QueryParser(&page); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	list, err := h.services.Winner.GetAllMonths(page)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(list)

}
