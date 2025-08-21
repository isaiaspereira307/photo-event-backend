package main

import (
	"log"

	"github.com/isaiaspereira307/photo-event-backend/configs"
	"github.com/isaiaspereira307/photo-event-backend/routes"
)

// @title           Photo Event API
// @version         1.0
// @description     API para gerenciamento de eventos fotográficos
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:9000
// @BasePath  /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Bearer token para autenticação
func main() {
	err := configs.Load()
	if err != nil {
		log.Fatal("Failed to load configurations:", err)
		panic(err)
	}

	err = configs.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	routes.Initialize()
}
