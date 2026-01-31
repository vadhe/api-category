package database

import (
	"database/sql"
	"fmt"
	"os"

	"log"
	"net/url"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func OpenPostgres() (*sql.DB, error) {
	err := godotenv.Load()
	if _, exists := os.LookupEnv("RAILWAY_ENVIRONMENT"); !exists {
		err := godotenv.Load()
		if err != nil {
			log.Println("Note: Error loading .env file locally (expected in Railway deployment)")
		}
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := url.QueryEscape(os.Getenv("DB_PASSWORD"))
	dbName := os.Getenv("DB_NAME")
	dbURL := os.Getenv("DB_URL")
	dbPort := os.Getenv("DB_PORT")
	dbSSLMode := os.Getenv("DB_SSL_MODE")

	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=%v", dbUser, dbPassword, dbURL, dbPort, dbName, dbSSLMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
