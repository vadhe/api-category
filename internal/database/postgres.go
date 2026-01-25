package database

import (
	"database/sql"
	"fmt"
	"os"

	"log"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func OpenPostgres() (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	connStr := fmt.Sprintf("postgres://%v:%v@localhost:5432/%v?sslmode=disable", dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
