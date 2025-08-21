package handlers

import (
	"github.com/isaiaspereira307/photo-event-backend/configs"
	users_handler "github.com/isaiaspereira307/photo-event-backend/handlers/users_handler"

	"gorm.io/gorm"
)

var db *gorm.DB

// @SecurityDefinitions.apiKey Bearer
// @in header
// @name Authorization
func InitializeHandlers() {
	db = configs.GetDB()

	users_handler.InitializeUserHandlers(db)
}
