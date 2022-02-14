package v1

import (
	"HundredToFive/internal/domain"
	"HundredToFive/pkg/logger"
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/api/option"
	"os"
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

func (h *Handler) firebaseNotification(ctx context.Context, noty domain.Notification, tokens []string, id int) (*messaging.BatchResponse, error) {

	opt := option.WithCredentialsFile(os.Getenv("FIREBASE_TOKEN"))

	config := &firebase.Config{ProjectID: os.Getenv("FIREBASE_PROJECT_ID")}

	app, err := firebase.NewApp(ctx, config, opt)

	if err != nil {
		return nil, fmt.Errorf("middlware.firebaseNotification: %w", err)
	}

	cl, err := app.Messaging(ctx)

	if err != nil {
		return nil, fmt.Errorf("middlware.firebaseNotification: %w", err)
	}

	message := &messaging.MulticastMessage{
		Data: map[string]string{
			"title":        noty.Title,
			"text":         noty.Text,
			"link":         noty.Link,
			"partner_logo": noty.Logo,
		},

		Tokens: tokens,
	}

	res, err := cl.SendMulticast(ctx, message)

	if err != nil {
		return nil, fmt.Errorf("middlware.firebaseNotification: %w", err)
	}
	input := domain.Notification{Status: 2}

	err = h.services.Notification.Update(id, input)

	if err != nil {
		return nil, fmt.Errorf("middlware.firebaseNotification: %w", err)
	}
	return res, nil
}
