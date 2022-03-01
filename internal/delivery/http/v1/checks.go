package v1

import (
	"HundredToFive/internal/domain"
	"HundredToFive/pkg/excel"
	"HundredToFive/pkg/media"
	"bytes"
	"errors"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"mime/multipart"
	"strconv"
	"time"
)

func (h *Handler) initCheckCategoryRoutes(api fiber.Router) {
	check := api.Group("/check")
	{
		check.Get("/download", h.downloadChecks)
		check.Get("", h.getAllChecks)
		check.Get("/:id", h.getCheckById)
		admin := check.Group("", jwtware.New(
			jwtware.Config{
				SigningKey: []byte(h.signingKey),
			}))
		{
			admin.Post("/register", h.registerCheck)
		}

	}
}

type CheckInfo struct {
	PartnerId   int                     `json:"partner_id" form:"partner_id" validate:"required"`
	CheckAmount int                     `json:"check_amount" form:"check_amount" validate:"required"`
	CheckDate   float64                 `json:"check_date" form:"check_date" validate:"required"`
	CheckImage  []*multipart.FileHeader `json:"check_image" form:"check_image" validate:"required" swaggerType:"array"`
}

// @Security User_Auth
// @Tags check
// @Description check upload
// @ModuleID registerCheck
// @Accept multipart/form-data
// @Produce json
// @Param partner_id formData int true "partner id"
// @Param check_amount formData int true "check cash"
// @Param check_date formData int true "date of check"
// @Param check_image formData file true "upload checks"
// @Success 201 {object} idResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /check/register [post]
func (h *Handler) registerCheck(c *fiber.Ctx) error {
	userType, _ := getUser(c)

	if userType != "user" {
		return c.Status(fiber.StatusUnauthorized).JSON(response{Message: "нет доступа"})
	}
	var input CheckInfo

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	_, userId := getUser(c)
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	files := form.File["check_image"]

	fileList := make([]string, 0)

	var image string
	for _, file := range files {

		image, err = media.GetFileName(c, file)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}

		fileList = append(fileList, image)

	}

	checks := domain.CheckInfo{
		UserId:      userId,
		PartnerId:   input.PartnerId,
		CheckAmount: input.CheckAmount,
		CheckDate:   input.CheckDate,
		CheckList:   fileList,
	}

	id, err := h.services.Check.Create(checks)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(idResponse{ID: id})
}

// @Tags check
// @Description gets all checks
// @ID get-all-checks
// @Accept json
// @Produce json
// @Param filter query  domain.FilterForCheck true "filter by money amount and partner"
// @Param array query domain.Pagination  true "A page info"
// @Success 200 {object} domain.GetAllChecksCategoryResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /check [get]
func (h *Handler) getAllChecks(c *fiber.Ctx) error {
	var (
		filterForCHeck domain.FilterForCheck
		err            error
	)

	if err := c.QueryParser(&filterForCHeck); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	var page domain.Pagination
	if err := c.QueryParser(&page); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	list, err := h.services.Check.GetAll(c, page, filterForCHeck)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(list)
}

// @Tags check
// @Description download members to excel
// @ID downloadChecks
// @Produce application/octet-stream
// @Param filter query  domain.FilterForCheck true "filter by money amount and partner"
// @Success 200 {file} file binary
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /check/download [get]
func (h *Handler) downloadChecks(c *fiber.Ctx) error {

	var filterForCheck domain.FilterForCheck

	if err := c.QueryParser(&filterForCheck); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	cellValue := map[string]string{"A1": "id", "B1": "id пользователя", "C1": "имя", "D1": "номер телефона", "E1": "партнер", "F1": "дата загрузки чека", "G1": "дата чека", "H1": "сумма чека", "I1": "изображение чека"}

	file, err := excel.File(cellValue)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}

	file, err = h.services.Check.DownloadChecks(c, file, filterForCheck)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}

	var b bytes.Buffer

	if err := file.Write(&b); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	downloadName := time.Now().UTC().Format("checks.xlsx")
	c.GetRespHeader("Content-Description", "File Transfer")
	c.GetRespHeader("Content-Type", "application/octet-stream")
	c.GetRespHeader("Content-Disposition", "attachment; filename="+downloadName)

	err = file.SaveAs(downloadName)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}

	return c.Download(downloadName)

}

// @Tags check
// @Description get check by id
// @ID get-check-by-id
// @Accept json
// @Produce json
// @Param id path string true "check id"
// @Success 200 {object} domain.UserChecks
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /check/{id} [get]
func (h *Handler) getCheckById(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	list, err := h.services.Check.GetById(c, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(list)
}
