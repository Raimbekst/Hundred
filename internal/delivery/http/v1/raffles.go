package v1

import (
	"HundredToFive/internal/domain"
	"HundredToFive/pkg/excel"
	"HundredToFive/pkg/validation/validationStructs"
	"bytes"
	"errors"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"strconv"
	"time"
)

func (h *Handler) initRaffleCategoryRoutes(api fiber.Router) {
	raffle := api.Group("/raffle")
	{
		raffle.Get("/download", h.downloadRaffles)
		raffle.Get("/", h.getAllRaffles)
		raffle.Get("/:id", h.getRaffleById)
		admin := raffle.Group("", jwtware.New(
			jwtware.Config{
				SigningKey: []byte(h.signingKey),
			}))
		{
			admin.Post("", h.createRaffle)
			admin.Put("/:id", h.updateRaffle)
			admin.Delete("/:id", h.deleteRaffle)
		}
	}
}

type Raffle struct {
	RaffleDate    float64 `json:"raffle_date" validate:"required"`
	RaffleTime    int     `json:"raffle_time" validate:"required"`
	CheckCategory int     `json:"check_category" validate:"required"`
	RaffleType    int     `json:"raffle_type" validate:"required" enums:"1,2,3" example:"1" default:"1"`
	Reference     string  `json:"reference"`
}

// @Security User_Auth
// @Tags raffle
// @ModuleID createRaffle
// @Accept  json
// @Produce  json
// @Param data body Raffle  true "game"
// @Success 201 {object} idResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /raffle [post]
func (h *Handler) createRaffle(c *fiber.Ctx) error {

	userType, _ := getUser(c)

	if userType != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(response{Message: "нет доступа"})
	}

	var input Raffle
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	ok, errs := validationStructs.ValidateStruct(input)

	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(errs)
	}

	raffle := domain.Raffle{

		RaffleDate:    input.RaffleDate,
		RaffleTime:    input.RaffleTime,
		CheckCategory: input.CheckCategory,
		RaffleType:    input.RaffleType,
		Reference:     input.Reference,
	}

	id, err := h.services.Raffle.Create(raffle)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(idResponse{ID: id})
}

// @Tags raffle
// @Description get all raffles
// @ID get-all-raffles
// @Accept  json
// @Produce  json
// @Param filter query domain.FilterForRaffles true "A filter info"
// @Param array query domain.Pagination  true "A page info"
// @Success 200 {object} domain.GetAllRaffleCategoryResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /raffle [get]
func (h *Handler) getAllRaffles(c *fiber.Ctx) error {
	var filter domain.FilterForRaffles

	if err := c.QueryParser(&filter); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	var page domain.Pagination
	if err := c.QueryParser(&page); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	list, err := h.services.Raffle.GetAll(page, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(list)

}

// @Tags raffle
// @Description download raffles to excel
// @ID downloadRaffles
// @Produce application/octet-stream
// @Success 200 {file} file binary
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /raffle/download [get]
func (h *Handler) downloadRaffles(c *fiber.Ctx) error {
	cellValue := map[string]string{"A1": "id", "B1": "имя пбедителя", "C1": "телефон номера", "D1": "дата розыгрыша", "E1": "категрия чека", "F1": "тип розыгрыша", "G1": "статус"}

	file, err := excel.File(cellValue)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}

	file, err = h.services.Raffle.DownloadRaffles(file)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}

	var b bytes.Buffer

	if err := file.Write(&b); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	downloadName := time.Now().UTC().Format("raffles.xlsx")
	c.GetRespHeader("Content-Description", "File Transfer")
	c.GetRespHeader("Content-Type", "application/octet-stream")
	c.GetRespHeader("Content-Disposition", "attachment; filename="+downloadName)

	err = file.SaveAs(downloadName)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}

	return c.Download(downloadName)
}

// @Tags raffle
// @Description get raffle by id
// @ID get-raffle-by-id
// @Accept  json
// @Produce  json
// @Param id path string true "raffle id"
// @Success 200 {object} domain.Raffle
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /raffle/{id} [get]
func (h *Handler) getRaffleById(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	list, err := h.services.Raffle.GetById(id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(list)
}

type UpdateRaffle struct {
	RaffleDate    float64 `json:"raffle_date" `
	RaffleTime    int     `json:"raffle_time"`
	CheckCategory int     `json:"check_category" `
	RaffleType    int     `json:"raffle_type"  enums:"1,2,3" example:"1" default:"1"`
	Reference     string  `json:"reference"`
}

// @Security User_Auth
// @Tags raffle
// @Description  update  raffle
// @ModuleID updateRaffle
// @Accept  json
// @Produce  json
// @Param id path string true "raffle id"
// @Param input body Raffle false "raffle"
// @Success 200 {object} okResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /raffle/{id} [put]
func (h *Handler) updateRaffle(c *fiber.Ctx) error {
	userType, _ := getUser(c)

	if userType != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(response{Message: "нет доступа"})
	}

	var input UpdateRaffle
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	raffle := domain.Raffle{

		RaffleDate:    input.RaffleDate,
		RaffleTime:    input.RaffleTime,
		CheckCategory: input.CheckCategory,
		RaffleType:    input.RaffleType,
		Reference:     input.Reference,
	}

	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	if err := h.services.Raffle.Update(id, raffle); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(okResponse{Message: "OK"})

}

// @Security User_Auth
// @Tags raffle
// @Description delete raffle
// @ModuleID deleteRaffle
// @Accept  json
// @Produce  json
// @Param id path string true "raffle id"
// @Success 200 {object} okResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /raffle/{id} [delete]
func (h *Handler) deleteRaffle(c *fiber.Ctx) error {
	userType, _ := getUser(c)

	if userType != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(response{Message: "нет доступа"})
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	if err := h.services.Raffle.Delete(id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(okResponse{Message: "OK"})
}
