package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/isaiaspereira307/photo-event-backend/handlers"
	auth_2fa_handler "github.com/isaiaspereira307/photo-event-backend/handlers/auth_2fa_handler"
)

func InitializeAuthRoutes(router *gin.Engine) {
	router.POST("/api/v1/auth/login", handlers.Login)
	router.POST("/api/v1/auth/register", handlers.Register)
	router.POST("/api/v1/auth/logout", handlers.Logout)
	router.POST("/api/v1/auth/refresh-token", handlers.RefreshToken)
	router.POST("/api/v1/auth/forgot-password", handlers.ForgotPassword)
	router.POST("/api/v1/auth/reset-password", handlers.ResetPassword)

	router.POST("/api/v1/auth/2fa/verify", auth_2fa_handler.VerifyLoginCode)

	// Rotas protegidas de 2FA (requerem autenticação)
	twoFaRoutes := router.Group("/api/v1/auth/2fa")
	{
		twoFaRoutes.GET("/settings", auth_2fa_handler.Get2FASettings)
		twoFaRoutes.POST("/settings", auth_2fa_handler.EnableDisable2FA)
	}
}
