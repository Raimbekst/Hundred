package app

import (
	"HundredToFive/internal/config"
	delivery "HundredToFive/internal/delivery/http"
	repos "HundredToFive/internal/repository"
	"HundredToFive/internal/server"
	"HundredToFive/internal/service"
	"HundredToFive/pkg/auth"
	"HundredToFive/pkg/database"
	"HundredToFive/pkg/database/redis"
	"HundredToFive/pkg/email/smtp"
	"HundredToFive/pkg/hash"
	"HundredToFive/pkg/logger"
	"HundredToFive/pkg/phone"
	"context"
	"errors"
	"github.com/robfig/cron/v3"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(configPath string) {
	cfg, err := config.Init(configPath)
	if err != nil {
		logger.Error(err)
	}
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		logger.Error(err)
	}
	red, err := redis.NewRedisDB(cfg)
	if err != nil {
		logger.Error(err)
	}
	emailSender, err := smtp.NewSMTPSender(cfg.SMTP.From, cfg.SMTP.Pass, cfg.SMTP.Host, cfg.SMTP.Port)
	if err != nil {
		logger.Error(err)

		return
	}
	hashes := hash.NewSHA1Hashes(cfg.Auth.PasswordSalt)

	otpNumberGenerator := phone.NewSecretGenerator()

	tokenManager, err := auth.NewManager(cfg.Auth.JWT.SigningKey)
	if err != nil {
		logger.Error(err)
	}

	cron := cron.New()

	ctx := context.TODO()

	repository := repos.NewRepository(db)
	services := service.NewService(service.Deps{
		Repos:           repository,
		Hashes:          hashes,
		OtpPhone:        otpNumberGenerator,
		Ctx:             ctx,
		Redis:           red,
		TokenManager:    tokenManager,
		AccessTokenTTL:  cfg.Auth.JWT.AccessTokenTTL,
		RefreshTokenTTL: cfg.Auth.JWT.RefreshTokenTTL,
		EmailSender:     emailSender,
		EmailConfig:     cfg.Email,
	})

	handlers := delivery.NewHandler(services, tokenManager, cfg.Auth.JWT.SigningKey, ctx, cron)

	srv := server.NewServer(handlers.Init(cfg))

	go func() {
		if err := srv.Run(cfg); !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("error occur while running http server: %s\n", err.Error())
		}
	}()

	logger.Info("server started")

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 3 * time.Second

	_, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(); err != nil {
		logger.Errorf("failed to stop server: %v", err)
	}
}
