package v1

import (
	"HundredToFive/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/golang-jwt/jwt/v4"
	"log"
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

type client struct{} // Add more data to this type if needed

var (
	clients    = make(map[*websocket.Conn]client) // Note: although large maps with pointer-like types (e.g. strings) as keys are slow, using pointers themselves as keys is acceptable and fast
	register   = make(chan *websocket.Conn)
	broadcast  = make(chan string)
	unregister = make(chan *websocket.Conn)
)

func runHub(interface{}) {
	for {
		select {
		case connection := <-register:
			clients[connection] = client{}
			log.Println("connection registered")

		case message := <-broadcast:
			log.Println("message received:", message)

			// Send the message to all clients
			for connection := range clients {
				if err := connection.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
					log.Println("write error:", err)

					connection.WriteMessage(websocket.CloseMessage, []byte{})
					connection.Close()
					delete(clients, connection)
				}
			}

		case connection := <-unregister:
			// Remove the client from the hub
			delete(clients, connection)

			log.Println("connection unregistered")
		}
	}
}
