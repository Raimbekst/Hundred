package v1

import (
	"HundredToFive/internal/service"
	"HundredToFive/pkg/auth"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	services     *service.Service
	tokenManager auth.TokenManager
	signingKey   string
}

func NewHandler(services *service.Service, tokenManager auth.TokenManager, signingKey string) *Handler {
	return &Handler{services: services, tokenManager: tokenManager, signingKey: signingKey}
}

func (h *Handler) Init(api fiber.Router) {
	v1 := api.Group("/v1")
	{
		h.initCheckCategoryRoutes(v1)
		h.initWinnerCategoryRoutes(v1)
		h.initPartnerCategoryRoutes(v1)
		h.initUserRoutes(v1)
		h.initCityCategoryRoutes(v1)
		h.initRaffleCategoryRoutes(v1)
		h.initBannerCategoryRoutes(v1)
		h.initDescCategoryRoutes(v1)
		h.initAboutWebsiteRoutes(v1)
		h.initFaqCategoryRoutes(v1)
		h.initNotificationRoutes(v1)
		h.initConditionRoutes(v1)

	}
}
