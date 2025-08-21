package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isaiaspereira307/photo-event-backend/internal/models"
	"github.com/isaiaspereira307/photo-event-backend/schemas"
)

// @BasePath /api/v1
// @Summary Update an user
// @Description Update an user
// @Tags user
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "User ID"
// @Param request body UpdateUserRequest true "Update User Request"
// @Success 200 {object} UpdateUserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id} [put]
func UpdateUser(ctx *gin.Context) {
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

	var req UpdateUserRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Message: "id is required"})
		return
	}

	new := schemas.User{}
	if err := db.First(&new, id).Error; err != nil {
		logger.Errorf("error: %s", err.Error())
		sendErr(ctx, err, http.StatusInternalServerError)
		return
	}

	if req.FirstName != "" && req.LastName != "" {
		new.Nome = req.FirstName + " " + req.LastName
	} else if req.FirstName != "" {
		new.Nome = req.FirstName
	}
	if req.Email != "" {
		new.Email = req.Email
	}
	if req.Password != "" {
		hashedPassword, err := hashPassword(req.Password)
		if err != nil {
			sendErr(ctx, err, http.StatusInternalServerError)
			return
		}
		new.Senha = hashedPassword
	}
	// Os campos IsActive, IsStaff, IsSuperuser, IsPremium não existem mais no esquema atual
	// Se necessário, podemos mapear para o campo Tipo (UserRole)

	if err := db.Save(&new).Error; err != nil {
		logger.Errorf("error updating opening: %v", err.Error())
		sendErr(ctx, err, http.StatusInternalServerError)
		return
	}

	sendSuccess(ctx, "updateUser", new)
}
