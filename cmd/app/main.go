package main

import "HundredToFive/internal/app"

const (
	configPath = "configs"
)

// @title A Hundred to Five
// @version 2.0
// @description API Server for  A Hundred to Five

// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1/

// @securityDefinitions.apikey User_Auth
// @in header
// @name Authorization
func main() {
	app.Run(configPath)

}
