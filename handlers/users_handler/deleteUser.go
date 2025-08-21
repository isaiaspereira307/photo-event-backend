package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isaiaspereira307/photo-event-backend/internal/models"
	"github.com/isaiaspereira307/photo-event-backend/schemas"
)

// @BasePath /api/v1
// @Summary Delete an user
// @Description Delete an user
// @Tags user
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Delete User Param"
// @Success 200 {object} DeleteUserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /users/{id} [delete]
func DeleteUser(ctx *gin.Context) {
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

	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	user := schemas.User{}
	if err := db.First(&user, id).Error; err != nil {
		sendErr(ctx, err, http.StatusNotFound)
		return
	}

	err := db.Delete(&user).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
		return
	}

	sendSuccess(ctx, "deleteUser", user)
}
