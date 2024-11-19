package database

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is the global database connection
var DB *gorm.DB

// ConnectDatabase initializes the database connection
func ConnectDatabase() {
	dsn := os.Getenv("DB_CONNECTION")
	if dsn == "" {
		log.Fatal("DB_CONNECTION environment variable not set")
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate models
	err = DB.AutoMigrate(&User{}, &Conversation{}, &Message{})
	if err != nil {
		log.Fatalf("Failed to auto-migrate models: %v", err)
	}

	log.Println("Database connected successfully")
}
