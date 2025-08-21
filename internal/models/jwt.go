package models

import "github.com/golang-jwt/jwt/v5"

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claims struct {
	Email       string `json:"email"`
	IsActive    bool   `json:"is_active"`
	IsStaff     bool   `json:"is_staff"`
	IsSuperuser bool   `json:"is_superuser"`
	IsPremium   bool   `json:"is_premium"`
	jwt.RegisteredClaims
}
