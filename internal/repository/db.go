package repository

import (
	"fmt"
	"log"
	"os"

	"github.com/ankitsingh/urlshortener/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("Connected to Database")

	// Auto migrate models
	err = DB.AutoMigrate(&models.User{}, &models.Link{}, &models.ClickEvent{})
	if err != nil {
		log.Fatalf("Failed to auto migrate database: %v", err)
	}
}
