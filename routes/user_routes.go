package routes

import (
	"github.com/gin-gonic/gin"
	users_handlers "github.com/isaiaspereira307/photo-event-backend/handlers/users_handler"
	"github.com/isaiaspereira307/photo-event-backend/middleware"
)

func InitializeUserRoutes(router *gin.RouterGroup) {
	protected := router.Group("/")
	protected.Use(middleware.RequireStaff())
	{
		protected.GET("/users", users_handlers.GetUsers)
		protected.GET("/users/:id", users_handlers.GetUser)
		protected.POST("/users", users_handlers.CreateUser)
		protected.PUT("/users/:id", users_handlers.UpdateUser)
		protected.DELETE("/users/:id", users_handlers.DeleteUser)
	}
	router.GET("/auth/me", users_handlers.GetCurrentUser)
}
