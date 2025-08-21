package handlers

import (
	"time"

	"github.com/isaiaspereira307/photo-event-backend/schemas"
)

type ErrorResponse struct {
	Message   string `json:"message"`
	ErrorCode string `json:"error_code"`
}

// LoginResponse é o tipo retornado pelo login
type LoginResponse struct {
	Success     bool         `json:"success"`
	Message     string       `json:"message"`
	Requires2FA bool         `json:"requires_2fa"`
	Token       string       `json:"token,omitempty"`
	User        UserResponse `json:"user"`
}

// UserResponse contém informações básicas do usuário
type UserResponse struct {
	ID          uint      `json:"id"`
	Email       string    `json:"email"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	IsActive    bool      `json:"is_active,omitempty"`
	IsStaff     bool      `json:"is_staff,omitempty"`
	IsSuperuser bool      `json:"is_superuser,omitempty"`
	IsPremium   bool      `json:"is_premium,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

type LogoutResponse struct {
	Message string `json:"message"`
}

type RefreshTokenResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

type AuthResponse struct {
	User  schemas.User `json:"user"`
	Token string       `json:"token"`
}
