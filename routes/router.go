package routes

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/isaiaspereira307/photo-event-backend/configs"
)

var (
	logger *configs.Logger
)

func Initialize() {
	logger = configs.GetLogger("routes")
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60,
	}))

	InitializeRoutes(router)
	port := fmt.Sprintf(":%s", configs.GetServerPort())

	if err := router.Run(port); err != nil {
		logger.Errorf("Failed to run server: %v", err)
	}
}
