package handlers

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/isaiaspereira307/photo-event-backend/internal/models"
	"github.com/isaiaspereira307/photo-event-backend/schemas"
	"gorm.io/gorm"
)

var db *gorm.DB

func Initialize2FAHandlers(database *gorm.DB) {
	db = database
}

// Request/Response structures
type Enable2FARequest struct {
	Enable bool `json:"enable" binding:"required"`
}

type Enable2FAResponse struct {
	Success bool                               `json:"success"`
	Message string                             `json:"message"`
	Data    *schemas.TwoFactorSettingsResponse `json:"data,omitempty"`
}

type Verify2FARequest struct {
	Token string `json:"token" binding:"required"`
}

type Verify2FAResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
}

type RequestLoginResponse struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	Requires2FA bool   `json:"requires_2fa"`
	Token       string `json:"token,omitempty"`
}

type ErrorResponse struct {
	Message   string `json:"message"`
	ErrorCode string `json:"error_code"`
}

// EnableDisable2FA godoc
// @BasePath /api/v1
// @Summary Enable or disable 2FA
// @Description Enable or disable two-factor authentication for the current user
// @Tags auth
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body Enable2FARequest true "Enable 2FA Request"
// @Success 200 {object} Enable2FAResponse
// @Summary Habilitar/Desabilitar 2FA
// @Description Habilita ou desabilita a autenticação de dois fatores para o usuário
// @Tags auth-2fa
// @Security Bearer
// @Accept json
// @Produce json
// @Param request body Enable2FARequest true "Configuração 2FA"
// @Success 200 {object} Enable2FAResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /auth/2fa/settings [post]
func EnableDisable2FA(c *gin.Context) {
	claimsInterface, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token não encontrado"})
		return
	}

	claims, ok := claimsInterface.(*models.Claims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar token"})
		return
	}

	var req Enable2FARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Buscar usuário
	var user schemas.User
	if err := db.Where("email = ?", claims.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		return
	}

	// Buscar ou criar configurações 2FA
	var settings schemas.TwoFactorSettings
	result := db.Where("user_id = ?", user.ID).First(&settings)
	if result.Error != nil {
		// Criar configurações se não existirem
		settings = schemas.TwoFactorSettings{
			UserID:  user.ID,
			Enabled: req.Enable,
		}
		if err := db.Create(&settings).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar configurações 2FA"})
			return
		}
	} else {
		// Atualizar configurações existentes
		settings.Enabled = req.Enable
		if err := db.Save(&settings).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar configurações 2FA"})
			return
		}
	}

	// Preparar resposta
	settingsResponse := schemas.TwoFactorSettingsResponse{
		ID:             settings.ID,
		UserID:         settings.UserID,
		Enabled:        settings.Enabled,
		LastVerifiedAt: settings.LastVerifiedAt,
		CreatedAt:      settings.CreatedAt,
		UpdatedAt:      settings.UpdatedAt,
	}

	c.JSON(http.StatusOK, Enable2FAResponse{
		Success: true,
		Message: fmt.Sprintf("2FA %s com sucesso", map[bool]string{true: "ativado", false: "desativado"}[req.Enable]),
		Data:    &settingsResponse,
	})
}

// Get2FASettings godoc
// @BasePath /api/v1
// @Summary Get 2FA settings
// @Summary Obter configurações 2FA
// @Description Retorna as configurações atuais de 2FA do usuário
// @Tags auth-2fa
// @Security Bearer
// @Produce json
// @Success 200 {object} Enable2FAResponse
// @Failure 401 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /auth/2fa/settings [get]
func Get2FASettings(c *gin.Context) {
	claimsInterface, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token não encontrado"})
		return
	}

	claims, ok := claimsInterface.(*models.Claims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar token"})
		return
	}

	// Buscar usuário
	var user schemas.User
	if err := db.Where("email = ?", claims.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		return
	}

	// Buscar configurações 2FA
	var settings schemas.TwoFactorSettings
	result := db.Where("user_id = ?", user.ID).First(&settings)
	if result.Error != nil {
		// Retornar configurações padrão se não existirem
		c.JSON(http.StatusOK, Enable2FAResponse{
			Success: true,
			Message: "Configurações 2FA",
			Data: &schemas.TwoFactorSettingsResponse{
				UserID:  user.ID,
				Enabled: false,
			},
		})
		return
	}

	// Preparar resposta
	settingsResponse := schemas.TwoFactorSettingsResponse{
		ID:             settings.ID,
		UserID:         settings.UserID,
		Enabled:        settings.Enabled,
		LastVerifiedAt: settings.LastVerifiedAt,
		CreatedAt:      settings.CreatedAt,
		UpdatedAt:      settings.UpdatedAt,
	}

	c.JSON(http.StatusOK, Enable2FAResponse{
		Success: true,
		Message: "Configurações 2FA",
		Data:    &settingsResponse,
	})
}

// Generate2FAToken gera um código de 6 dígitos para verificação 2FA
func Generate2FAToken(userID uuid.UUID) (string, error) {
	// Limpar tokens antigos
	db.Where("user_id = ? OR expires_at < ?", userID, time.Now()).Delete(&schemas.TwoFactorToken{})

	// Gerar novo token de 6 dígitos
	rand.Seed(time.Now().UnixNano())
	token := fmt.Sprintf("%06d", rand.Intn(1000000))

	// Salvar token no banco
	expiresAt := time.Now().Add(10 * time.Minute)
	newToken := schemas.TwoFactorToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
	}

	if err := db.Create(&newToken).Error; err != nil {
		return "", err
	}

	return token, nil
}

// VerifyLoginCode godoc
// @BasePath /api/v1
// @Summary Verify 2FA code
// @Description Verify the 2FA verification code during login
// @Tags auth
// @Accept json
// @Produce json
// @Param email query string true "User email"
// @Param request body Verify2FARequest true "Verification Code"
// @Success 200 {object} Verify2FAResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/2fa/verify [post]
func VerifyLoginCode(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email é obrigatório"})
		return
	}

	var req Verify2FARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Buscar usuário
	var user schemas.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		return
	}

	// Verificar token
	var token schemas.TwoFactorToken
	if err := db.Where("user_id = ? AND token = ? AND expires_at > ? AND used = ?",
		user.ID, req.Token, time.Now(), false).First(&token).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Código inválido ou expirado"})
		return
	}

	// Marcar token como usado
	token.Used = true
	if err := db.Save(&token).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar verificação"})
		return
	}

	// Atualizar última verificação
	var settings schemas.TwoFactorSettings
	if err := db.Where("user_id = ?", user.ID).First(&settings).Error; err == nil {
		settings.LastVerifiedAt = time.Now()
		db.Save(&settings)
	}

	// Gerar token JWT
	jwtToken, err := models.GenerateToken(user.Email, 24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar token"})
		return
	}

	c.JSON(http.StatusOK, Verify2FAResponse{
		Success: true,
		Message: "Verificação concluída com sucesso",
		Token:   jwtToken,
	})
}
