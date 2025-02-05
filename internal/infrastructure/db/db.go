package db

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// DB is a global variable that holds the database connection pool.
var DB *sqlx.DB

// InitDB initializes the database connection using environment variables.
func InitDB() (*sqlx.DB, error){

	// Load environment variables from the .env file.
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Retrieve the DATABASE_URL from the environment variables.
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	// Connect to the PostgreSQL database using the provided DSN.
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Assign the connected database to the global DB variable.
	DB = db
	return db, nil
}