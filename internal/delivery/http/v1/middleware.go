package v1

import (
	"context"
	"errors"
	"fmt"
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

func (h *Handler) changeStatus() {
	_, _ = h.cron.AddFunc("* * * * *", func() {
		timeNow := time.Now().Unix()
		err := h.services.Raffle.UpdateStatus(timeNow)

		if err != nil {

			fmt.Println(err)
		}
		fmt.Println("changed")
	})
}

func (h *Handler) sendNotification(ctx context.Context) {

	_, _ = h.cron.AddFunc("* * * * *", func() {

		timeNow := time.Now().Unix()

		list, err := h.services.Notification.GetNotificationByDate(timeNow)

		if err != nil {
			fmt.Println(err)
		}

		tokens, err := h.services.Notification.GetAllRegistrationTokens()

		if err != nil {
			fmt.Println(err)

		}
		if tokens == nil {
			tokens = []string{"Random_token"}
		}
		if list != nil {

			for _, value := range list {

				res, err := h.firebaseNotification(ctx, value, tokens, value.Id)

				if err != nil {
					fmt.Println(err)
				}

				fmt.Println(res.SuccessCount)

			}
		}

		fmt.Println("done")

	})

	h.cron.Start()

}
