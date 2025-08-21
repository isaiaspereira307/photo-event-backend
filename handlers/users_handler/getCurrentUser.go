package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isaiaspereira307/photo-event-backend/internal/models"
	"github.com/isaiaspereira307/photo-event-backend/schemas"
)

// @BasePath /api/v1
// @Summary Get current user
// @Description Get details of the currently authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} UserResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/me [get]
// @Summary Obter usuário atual
// @Description Retorna informações do usuário autenticado
// @Tags users
// @Security Bearer
// @Produce json
// @Success 200 {object} object{success=bool,message=string,data=schemas.UserResponse}
// @Failure 401 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /users/me [get]
func GetCurrentUser(ctx *gin.Context) {
	// Extract claims from context (set by auth middleware)
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

	// Get user from database using the email from claims
	var user schemas.User
	if err := db.Where("email = ?", claims.Email).First(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar usuário"})
		return
	}

	userResponse := schemas.UserResponse{
		ID:           user.ID,
		Nome:         user.Nome,
		Email:        user.Email,
		Tipo:         user.Tipo,
		Creditos:     user.Creditos,
		DataCadastro: user.DataCadastro,
		UltimoAcesso: user.UltimoAcesso,
		Biografia:    user.Biografia,
		URLPerfil:    user.URLPerfil,
		URLCapa:      user.URLCapa,
		Telefone:     user.Telefone,
		Website:      user.Website,
	}

	// Return user data
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User data retrieved successfully",
		"data":    userResponse,
	})
}
