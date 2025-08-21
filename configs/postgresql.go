package configs

import (
	"fmt"

	"github.com/isaiaspereira307/photo-event-backend/schemas"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitializePostgres(cfg DBConfig) (*gorm.DB, error) {
	logger := GetLogger("postgres")
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.Database, cfg.Port, cfg.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Errorf("Error connecting to PostgreSQL: %s", err.Error())
		return nil, err
	}

	logger.Info("Connected to PostgreSQL successfully")

	err = db.AutoMigrate(
		&schemas.User{},
		&schemas.TwoFactorSettings{},
		&schemas.TwoFactorToken{},
	)
	if err != nil {
		logger.Errorf("Error migrating PostgreSQL: %s", err.Error())
		return nil, err
	}

	return db, nil
}
