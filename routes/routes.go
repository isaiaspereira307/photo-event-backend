package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/isaiaspereira307/photo-event-backend/configs"
	"github.com/isaiaspereira307/photo-event-backend/docs"
	"github.com/isaiaspereira307/photo-event-backend/handlers"
	"github.com/isaiaspereira307/photo-event-backend/middleware"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitializeRoutes(router *gin.Engine) {
	db = configs.GetDB()
	url := configs.GetServerUrl()
	handlers.InitializeHandlers()
	basePath := "/api/v1"
	docs.SwaggerInfo.Title = "API Documentation"
	docs.SwaggerInfo.Description = "This is a sample server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = url
	docs.SwaggerInfo.BasePath = basePath
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	InitializeAuthRoutes(router)
	protected := router.Group(basePath)
	protected.Use(middleware.AuthMiddleware())
	{
		InitializeUserRoutes(protected)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
