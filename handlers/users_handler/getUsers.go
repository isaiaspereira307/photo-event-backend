package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isaiaspereira307/photo-event-backend/internal/models"
	"github.com/isaiaspereira307/photo-event-backend/schemas"
)

// @BasePath /api/v1
// @Summary Show users
// @Description Show all users
// @Tags user
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} ListUserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /users [get]
func GetUsers(ctx *gin.Context) {
	claimsInterface, exists := ctx.Get("claims")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token não encontrado"})
		return
	}

	claims, ok := claimsInterface.(*models.Claims)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar token"})
		return
	}

	// Check if user has permission
	if !claims.IsStaff {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Usuário sem permissão"})
		return
	}

	user := []schemas.User{}

	if err := db.Find(&user).Error; err != nil {
		sendErr(ctx, err, http.StatusInternalServerError)
		return
	}

	sendSuccess(ctx, "getUsers", user)
}
