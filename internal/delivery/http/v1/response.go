package v1

import (
	"HundredToFive/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"strconv"
)

const (
	user  = "user"
	admin = "admin"
)

type idResponse struct {
	ID interface{} `json:"id"`
}
type okResponse struct {
	Message string `json:"message"`
}

type response struct {
	Message string `json:"detail"`
}

func newResponse(c *fiber.Ctx, statusCode, message string) {
	logger.Error(message)

}

func getUser(c *fiber.Ctx) (string, int) {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id := claims["jti"].(string)
	userType := claims["sub"].(string)
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return "", 0
	}
	return userType, idInt
}
