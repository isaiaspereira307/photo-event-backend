package handlers

import "gorm.io/gorm"

var db *gorm.DB

func InitializeUserHandlers(database *gorm.DB) {
	db = database
}
