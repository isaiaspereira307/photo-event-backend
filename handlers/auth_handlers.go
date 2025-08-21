package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/isaiaspereira307/photo-event-backend/configs"
	"github.com/isaiaspereira307/photo-event-backend/schemas"
	"golang.org/x/crypto/bcrypt"
)

// @Summary Registrar usuário
// @Description Registra um novo usuário no sistema
// @Tags auth
// @Accept json
// @Produce json
// @Param request body schemas.RegisterRequest true "Dados de registro"
// @Success 201 {object} schemas.SuccessResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 409 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /auth/register [post]
func Register(c *gin.Context) {
	var req schemas.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, schemas.ErrorResponse{
			Error:   "Dados inválidos",
			Message: err.Error(),
		})
		return
	}

	db := configs.GetDB()

	// Verificar se o email já existe
	var existingUser schemas.User
	if err := db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, schemas.ErrorResponse{
			Error:   "Email já cadastrado",
			Message: "Este email já está sendo usado por outro usuário",
		})
		return
	}

	// Validar tipo de usuário
	if req.Tipo == "" {
		req.Tipo = schemas.RoleUsuario
	}

	if req.Tipo != schemas.RoleUsuario && req.Tipo != schemas.RoleFotografo && req.Tipo != schemas.RoleAdmin {
		c.JSON(http.StatusBadRequest, schemas.ErrorResponse{
			Error:   "Tipo de usuário inválido",
			Message: "Tipo deve ser: usuario, fotografo ou admin",
		})
		return
	}

	// Hash da senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Senha), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Error: "Erro interno do servidor",
		})
		return
	}

	// Criar usuário
	user := schemas.User{
		Nome:  req.Nome,
		Email: req.Email,
		Senha: string(hashedPassword),
		Tipo:  req.Tipo,
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Error:   "Erro ao criar usuário",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, schemas.SuccessResponse{
		Message: "Usuário criado com sucesso",
		Data:    user,
	})
}

// @Summary Login do usuário
// @Description Realiza login de um usuário no sistema
// @Tags auth
// @Accept json
// @Produce json
// @Param request body schemas.LoginRequest true "Credenciais de login"
// @Success 200 {object} schemas.LoginResponse
// @Failure 400 {object} schemas.ErrorResponse
// @Failure 401 {object} schemas.ErrorResponse
// @Failure 500 {object} schemas.ErrorResponse
// @Router /auth/login [post]
func Login(c *gin.Context) {
	var req schemas.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, schemas.ErrorResponse{
			Error:   "Dados inválidos",
			Message: err.Error(),
		})
		return
	}

	db := configs.GetDB()

	// Buscar usuário pelo email
	var user schemas.User
	if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, schemas.ErrorResponse{
			Error:   "Credenciais inválidas",
			Message: "Email ou senha incorretos",
		})
		return
	}

	// Verificar senha
	if err := bcrypt.CompareHashAndPassword([]byte(user.Senha), []byte(req.Senha)); err != nil {
		c.JSON(http.StatusUnauthorized, schemas.ErrorResponse{
			Error:   "Credenciais inválidas",
			Message: "Email ou senha incorretos",
		})
		return
	}

	// Atualizar último acesso
	db.Model(&user).Update("ultimo_acesso", time.Now())

	// Gerar token JWT
	token, err := generateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, schemas.ErrorResponse{
			Error: "Erro ao gerar token",
		})
		return
	}

	c.JSON(http.StatusOK, schemas.LoginResponse{
		Token:   token,
		User:    user,
		Message: "Login realizado com sucesso",
	})
}

// @Summary Esqueceu a senha
// @Description Inicia o processo de recuperação de senha (não implementado)
// @Tags auth
// @Produce json
// @Success 501 {object} schemas.ErrorResponse
// @Router /auth/forgot-password [post]
func ForgotPassword(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, schemas.ErrorResponse{
		Error:   "Funcionalidade não implementada",
		Message: "Reset de senha será implementado em versão futura",
	})
}

// @Summary Reset de senha
// @Description Redefine a senha do usuário (não implementado)
// @Tags auth
// @Produce json
// @Success 501 {object} schemas.ErrorResponse
// @Router /auth/reset-password [post]
func ResetPassword(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, schemas.ErrorResponse{
		Error:   "Funcionalidade não implementada",
		Message: "Reset de senha será implementado em versão futura",
	})
}

// @Summary Logout
// @Description Realiza logout do usuário (não implementado)
// @Tags auth
// @Security Bearer
// @Produce json
// @Success 200 {object} schemas.SuccessResponse
// @Router /auth/logout [post]
func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, schemas.SuccessResponse{
		Message: "Logout realizado com sucesso",
	})
}

// @Summary Refresh Token
// @Description Renova o token de acesso (não implementado)
// @Tags auth
// @Security Bearer
// @Produce json
// @Success 501 {object} schemas.ErrorResponse
// @Router /auth/refresh [post]
func RefreshToken(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, schemas.ErrorResponse{
		Error:   "Funcionalidade não implementada",
		Message: "Refresh de token será implementado em versão futura",
	})
}

// generateJWT gera token JWT baseado no teste.go
func generateJWT(user schemas.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &schemas.Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Tipo,
	}

	// Adicionar claims padrão do JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": claims.UserID,
		"email":   claims.Email,
		"role":    claims.Role,
		"exp":     expirationTime.Unix(),
		"iat":     time.Now().Unix(),
	})

	tokenString, err := token.SignedString(configs.GetJwtSecret())
	return tokenString, err
}
