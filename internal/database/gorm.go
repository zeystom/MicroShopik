package database

import (
	"MicroShopik/configs"
	"MicroShopik/internal/domain"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(cfg *configs.Config) error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	err = db.AutoMigrate(
		&domain.User{},
		&domain.Product{},
		&domain.Role{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	DB = db
	log.Println("Database connected and migrated successfully")
	return nil
}

func GetDB() *gorm.DB {
	return DB
}
