package main

import (
	"log"

	"github.com/masahiro0000/BizMatch/internal/app"
)

func main() {
	// Launch the application
	if err := app.Run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}