package main

import (
	"log"
	"os"
)

func loadConfiguration() {
	// Load environment variables (e.g., from .env file)
	if err := os.Setenv("DB_CONNECTION", "postgres://user:password@localhost:5432/wasatext"); err != nil {
		log.Fatalf("Failed to set environment variables: %s", err)
	}

	log.Println("Configuration loaded successfully")
}
