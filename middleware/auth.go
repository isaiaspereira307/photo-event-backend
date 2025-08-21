package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/isaiaspereira307/photo-event-backend/configs"
	"github.com/isaiaspereira307/photo-event-backend/schemas"
)

// AuthMiddleware - Middleware de autenticação baseado no teste.go
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, schemas.ErrorResponse{
				Error: "Token de autorização requerido",
			})
			c.Abort()
			return
		}

		// Extrair token do header "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, schemas.ErrorResponse{
				Error: "Formato de token inválido",
			})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		// Verificar e parsear token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return configs.GetJwtSecret(), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, schemas.ErrorResponse{
				Error: "Token inválido ou expirado",
			})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, schemas.ErrorResponse{
				Error: "Claims inválidas no token",
			})
			c.Abort()
			return
		}

		// Verificar se usuário ainda existe
		var user schemas.User
		db := configs.GetDB()
		userID := claims["user_id"].(string)
		if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, schemas.ErrorResponse{
				Error: "Usuário não encontrado",
			})
			c.Abort()
			return
		}

		// Adicionar informações do usuário ao contexto
		c.Set("user_id", userID)
		c.Set("user_email", claims["email"].(string))
		c.Set("user_role", claims["role"].(string))
		c.Next()
	}
}

// AdminMiddleware - Middleware para verificar se o usuário é admin
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetString("user_role")
		if userRole != string(schemas.RoleAdmin) {
			c.JSON(http.StatusForbidden, schemas.ErrorResponse{
				Error:   "Acesso negado",
				Message: "Apenas administradores podem acessar este recurso",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireStaff - Compatibilidade com middleware existente
func RequireStaff() gin.HandlerFunc {
	return AdminMiddleware()
}
