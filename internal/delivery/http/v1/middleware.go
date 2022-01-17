package v1

import (
	"HundredToFive/internal/domain"
	"context"
	"errors"
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"strings"
	"time"
)

const (
	getterAll   = "все"
	getterUser  = "выборочно"
	planned     = "запланирован"
	sent        = "отправлено"
	archive     = "архивный"
	active      = "действующий"
	dayGame     = "Ежедневный розыгрыш"
	weeklyGame  = "Еженедельный розыгрыш"
	monthlyGame = "Ежемесячный розыгрыш"
)

func (h *Handler) userIdentity(token string) (int, error) {
	id, _, err := h.parseAuthHeader(token)
	if err != nil {
		return 0, fmt.Errorf("v1.userIdentity: %w", err)
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return 0, fmt.Errorf("v1.userIdentity: %w", errors.New("id is invalid type"))
	}
	return idInt, nil
}

func (h *Handler) parseAuthHeader(token string) (string, string, error) {

	return h.tokenManager.Parse(token)
}

func parseRequestHost(c *fiber.Ctx) string {
	refererHeader := c.Get("Referer")
	refererParts := strings.Split(refererHeader, "/")

	hostParts := strings.Split(refererParts[2], ":")

	return hostParts[0]
}

func (h *Handler) scheduleNotification(ctx context.Context, noty domain.Notification, tokens []string, id int) error {

	s := gocron.NewScheduler(time.Local)

	executionTime := int(int64(noty.Date) - time.Now().Unix())

	_, err := s.Every(1).Day().StartAt(time.Now().Add(time.Duration(executionTime) * time.Second)).Do(func() {

		notificationList, err := h.services.Notification.GetById(id)

		if err != nil {
			fmt.Println(err)
			return
		}

		res, err := h.firebaseNotification(ctx, noty, tokens, notificationList.Id)

		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(res.SuccessCount)

	})
	if err != nil {
		return err
	}

	s.StartAsync()
	return nil
}
