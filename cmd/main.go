package main

import (
	"log"

	"github.com/masahiro0000/BizMatch/internal/infrastructure/db"
	"github.com/masahiro0000/BizMatch/internal/infrastructure/router"
)

func main() {

	// Initialize the database connection.
	_, err := db.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	// Create a new router instance for handling HTTP requests.
	router := router.NewRouter()

	// Start the HTTP server on port 8080.
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}