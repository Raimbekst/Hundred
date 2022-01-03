package v1

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"strings"
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
