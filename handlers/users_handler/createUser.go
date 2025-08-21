package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isaiaspereira307/photo-event-backend/internal/models"
	"github.com/isaiaspereira307/photo-event-backend/schemas"
	"golang.org/x/crypto/bcrypt"
)

// @BasePath /api/v1
// @Summary Create an user
// @Description Create an user
// @Tags user
// @Accept json
// @Produce json
// @Param request body CreateUserRequest true "Create User Params"
// @Security Bearer
// @Success 200 {object} CreateUserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users [post]
func CreateUser(ctx *gin.Context) {
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

	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		sendErr(ctx, err, http.StatusBadRequest)
		return
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		sendErr(ctx, err, http.StatusInternalServerError)
		return
	}

	new := schemas.User{
		Nome:  req.FirstName + " " + req.LastName, // Combinar primeiro e último nome
		Email: req.Email,
		Senha: hashedPassword,
		Tipo:  schemas.RoleUsuario, // Valor padrão
	}

	err = db.Create(&new).Error
	if err != nil {
		sendErr(ctx, err, http.StatusInternalServerError)
		return
	}

	sendSuccess(ctx, "createUser", new)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
