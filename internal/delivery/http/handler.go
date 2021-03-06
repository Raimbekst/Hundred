package http

import (
	"HundredToFive/docs"
	"HundredToFive/internal/config"
	v1 "HundredToFive/internal/delivery/http/v1"
	"HundredToFive/internal/service"
	"HundredToFive/pkg/auth"
	"context"
	"fmt"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/robfig/cron/v3"
)

type Handler struct {
	services     *service.Service
	tokenManager auth.TokenManager
	signingKey   string
	ctx          context.Context
	cron         *cron.Cron
}

func NewHandler(services *service.Service, tokenManager auth.TokenManager, signingKey string, ctx context.Context, cron *cron.Cron) *Handler {
	return &Handler{services: services, tokenManager: tokenManager, signingKey: signingKey, ctx: ctx, cron: cron}
}

func (h *Handler) Init(cfg *config.Config) *fiber.App {
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port)

	if cfg.Environment == config.Prod {
		docs.SwaggerInfo.Host = fmt.Sprintf("%s", cfg.HTTP.Host)
	}
	router := fiber.New(fiber.Config{BodyLimit: 100 * 1024 * 1024})
	router.Use(logger.New())

	var ConfigDefault = cors.Config{
		Next:         nil,
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
	}

	router.Use(cors.New(ConfigDefault))

	router.Get("/swagger/*", swagger.Handler)

	h.initApi(router)
	router.Static("/media", "media")
	return router
}

func (h *Handler) initApi(router *fiber.App) {
	handler := v1.NewHandler(h.services, h.tokenManager, h.signingKey, h.ctx, h.cron)
	api := router.Group("/api")
	{
		handler.Init(api)
	}
}
