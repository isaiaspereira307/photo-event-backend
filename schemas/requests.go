package schemas

import "github.com/google/uuid"

// DTOs para requests/responses (baseado no teste.go)
type RegisterRequest struct {
	Nome  string   `json:"nome" binding:"required,min=2" example:"João Silva"`
	Email string   `json:"email" binding:"required,email" example:"joao@exemplo.com"`
	Senha string   `json:"senha" binding:"required,min=6" example:"senha123"`
	Tipo  UserRole `json:"tipo,omitempty" example:"usuario"`
}

type LoginRequest struct {
	Email string `json:"email" binding:"required,email" example:"joao@exemplo.com"`
	Senha string `json:"senha" binding:"required" example:"senha123"`
}

type LoginResponse struct {
	Token   string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User    User   `json:"user"`
	Message string `json:"message" example:"Login realizado com sucesso"`
}

type ChangePasswordRequest struct {
	SenhaAtual string `json:"senha_atual" binding:"required"`
	NovaSenha  string `json:"nova_senha" binding:"required,min=6"`
}

type UpdateProfileRequest struct {
	Nome      *string `json:"nome,omitempty"`
	Biografia *string `json:"biografia,omitempty"`
	URLPerfil *string `json:"url_perfil,omitempty"`
	URLCapa   *string `json:"url_capa,omitempty"`
	Telefone  *string `json:"telefone,omitempty"`
	Website   *string `json:"website,omitempty"`
}

// Responses padronizadas
type ErrorResponse struct {
	Error   string `json:"error" example:"Erro de validação"`
	Message string `json:"message,omitempty" example:"Campo obrigatório não informado"`
}

type SuccessResponse struct {
	Message string      `json:"message" example:"Operação realizada com sucesso"`
	Data    interface{} `json:"data,omitempty"`
}

// Claims para JWT
type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Role   UserRole  `json:"role"`
}
