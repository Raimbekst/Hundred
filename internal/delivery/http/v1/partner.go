package v1

import (
	"HundredToFive/internal/domain"
	"HundredToFive/pkg/excel"
	"HundredToFive/pkg/media"
	"bytes"
	"errors"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/google/uuid"
	"mime/multipart"
	"strconv"
	"time"
)

func (h *Handler) initPartnerCategoryRoutes(api fiber.Router) {
	partner := api.Group("/partner")
	{
		partner.Get("/", h.getAllPartners)
		partner.Get("/download", h.downloadPartners)
		partner.Get("/:id", h.getPartnerById)
		admin := partner.Group("", jwtware.New(
			jwtware.Config{
				SigningKey: []byte(h.signingKey),
			}))
		{
			admin.Post("", h.createPartner)
			admin.Put("/:id", h.updatePartner)
			admin.Delete("/:id", h.deletePartner)
		}
	}
}

type Partner struct {
	PartnerName      string                `json:"partner_name" form:"partner_name"`
	Position         int                   `json:"position" form:"position"`
	Logo             *multipart.FileHeader `json:"logo" form:"logo"`
	LinkWebsite      string                `json:"link_website" form:"link_website"`
	Banner           *multipart.FileHeader `json:"banner" form:"banner"`
	BannerKz         *multipart.FileHeader `json:"banner_kz" form:"banner_kz"`
	Status           int                   `form:"status"  enums:"1,2" default:"1"`
	StartPartnership string                `json:"start_partnership" form:"start_partnership"`
	EndPartnership   string                `json:"end_partnership" form:"end_partnership"`
	PartnerPackage   string                `json:"partner_package" form:"partner_package"`
	Reference        string                `json:"reference" form:"reference"`
}

// @Security User_Auth
// @Tags partner
// @ModuleID createPartner
// @Accept  multipart/form-data
// @Produce  json
// @Param partner_name formData string true "partner name"
// @Param logo  formData file false  "logo"
// @Param link_website formData string false "link_website"
// @Param banner formData file false "banner"
// @Param banner_kz formData file false "banner kz"
// @Param position formData int false "position"
// @Param status formData int true "only 1 or 2" Enums(1,2)
// @Param start_partnership formData string false "start of partnership"
// @Param end_partnership formData string false "end of partnership"
// @Param partner_package formData string false "partner package"
// @Param reference formData string false "reference"
// @Success 201 {object} idResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /partner [post]
func (h *Handler) createPartner(c *fiber.Ctx) error {
	var (
		input Partner
		err   error
	)
	userType, _ := getUser(c)
	if userType != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(response{Message: "нет доступа"})
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	var logo string

	file, _ := c.FormFile("logo")

	if file != nil {
		logo, err = media.GetFileName(c, file)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
		}
	}

	var banner string

	file1, _ := c.FormFile("banner")

	if file1 != nil {
		banner, err = media.GetFileName(c, file1)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
		}
	}

	var banner2 string
	file2, _ := c.FormFile("banner_kz")

	if file2 != nil {
		banner2, err = media.GetFileName(c, file2)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
		}
	}

	partner := domain.Partner{
		PartnerName:      input.PartnerName,
		Logo:             logo,
		LinkWebsite:      input.LinkWebsite,
		Banner:           banner,
		BannerKz:         banner2,
		Status:           input.Status,
		StartPartnership: input.StartPartnership,
		EndPartnership:   input.EndPartnership,
		PartnerPackage:   input.PartnerPackage,
		Reference:        input.Reference,
		Position:         input.Position,
	}

	id, err := h.services.Partner.Create(partner)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(idResponse{ID: id})
}

type FilterForPartner struct {
	Status int `json:"status" form:"status" query:"status" enums:"1,2"`
}

// @Tags partner
// @Description get all partners
// @ID get-all-partner
// @Accept  json
// @Produce  json
// @Param filter query FilterForPartner false "filter by status"
// @Param array query domain.Pagination  true "A page info"
// @Success 200 {object} domain.GetAllPartnersCategoryResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /partner [get]
func (h *Handler) getAllPartners(c *fiber.Ctx) error {
	url := c.BaseURL()
	var filter FilterForPartner

	if err := c.QueryParser(&filter); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	var page domain.Pagination
	if err := c.QueryParser(&page); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	list, err := h.services.Partner.GetAll(page, filter.Status)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}

	for _, value := range list.Data {
		if value.Logo != "" {
			value.Logo = url + "/" + "media/" + value.Logo

		}
		if value.Banner != "" {
			value.Banner = url + "/" + "media/" + value.Banner
		}
		if value.BannerKz != "" {
			value.BannerKz = url + "/" + "media/" + value.BannerKz
		}
	}
	return c.Status(fiber.StatusOK).JSON(list)

}

// @Tags partner
// @Description download members to excel
// @ID downloadPartners
// @Produce application/octet-stream
// @Success 200 {file} file binary
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /partner/download [get]
func (h *Handler) downloadPartners(c *fiber.Ctx) error {
	url := c.BaseURL()
	cellValue := map[string]string{"A1": "id", "B1": "имя партнера", "C1": "ссылка на сайт", "D1": "статус", "E1": "примечание", "F1": "лого партнера", "G1": "баннер", "H1": "баннер kz", "I1": "дата начала партнерства", "J1": "дата окончания партнерства"}

	file, err := excel.File(cellValue)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}

	file, err = h.services.Partner.DownloadPartners(file, url)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}

	var b bytes.Buffer

	if err := file.Write(&b); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	downloadName := time.Now().UTC().Format("partners.xlsx")
	c.GetRespHeader("Content-Description", "File Transfer")
	c.GetRespHeader("Content-Type", "application/octet-stream")
	c.GetRespHeader("Content-Disposition", "attachment; filename="+downloadName)

	err = file.SaveAs(downloadName)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}

	return c.Download(downloadName)
}

// @Tags partner
// @Description get partner by id
// @ID get-partner-by-id
// @Accept  json
// @Produce  json
// @Param id path string true "partner id"
// @Success 200 {object} domain.Partner
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /partner/{id} [get]
func (h *Handler) getPartnerById(c *fiber.Ctx) error {
	url := c.BaseURL()
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	list, err := h.services.Partner.GetById(id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	if list.Logo != "" {
		list.Logo = url + "/" + "media/" + list.Logo

	}
	if list.Banner != "" {
		list.Banner = url + "/" + "media/" + list.Banner
	}
	if list.BannerKz != "" {
		list.BannerKz = url + "/" + "media/" + list.BannerKz

	}
	return c.Status(fiber.StatusOK).JSON(list)
}

// @Security User_Auth
// @Tags partner
// @Description  update  partner
// @ModuleID updatePartner
// @Accept json
// @Produce  json
// @Param id path string true "partner id"
// @Param data body domain.UpdatePartner true "update partner"
// @Success 200 {object} okResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /partner/{id} [put]
func (h *Handler) updatePartner(c *fiber.Ctx) error {
	var (
		input domain.UpdatePartner
		err   error
	)

	userType, _ := getUser(c)

	if userType != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(response{Message: "нет доступа"})
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	var logo string

	if input.Logo != nil {
		if *input.Logo != "" {
			fil := uuid.New().String()

			logo, err = media.Base64ToImage(*input.Logo, fil)

			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
			}
			input.Logo = &logo

		}
	}
	var bannerRus string

	if input.Banner != nil {
		if *input.Banner != "" {
			fil := uuid.New().String()

			bannerRus, err = media.Base64ToImage(*input.Banner, fil)

			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
			}
		}

		input.Banner = &bannerRus
	}
	var bannerKz string

	if input.BannerKz != nil {
		if *input.BannerKz != "" {
			fil := uuid.New().String()

			bannerKz, err = media.Base64ToImage(*input.BannerKz, fil)

			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
			}
		}

		input.BannerKz = &bannerKz
	}
	if err := h.services.Partner.Update(id, input); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(okResponse{Message: "OK"})

}

// @Security User_Auth
// @Tags partner
// @Description delete partner
// @ModuleID deletePartner
// @Accept  json
// @Produce  json
// @Param id path string true "partner id"
// @Success 200 {object} okResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /partner/{id} [delete]
func (h *Handler) deletePartner(c *fiber.Ctx) error {

	userType, _ := getUser(c)

	if userType != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(response{Message: "нет доступа"})
	}
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	if err := h.services.Partner.Delete(id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(okResponse{Message: "OK"})
}
