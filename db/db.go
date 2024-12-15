package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func Connect() *sql.DB {
	connStr := os.Getenv("DB_CONNECTION_STRING")
	if connStr == "" {
		log.Fatal("DB_CONNECTION_STRING is not set")
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
	}

	log.Println("Successfully connected to the database!")
	return db
}

func MigrateToDB(db *sql.DB) {
	log.Println("Starting migrate to database...")

	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS Users (
            Login TEXT NOT NULL UNIQUE,
            Password TEXT NOT NULL,
            CreatedAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
        )
    `)
	if err != nil {
		log.Printf("Error creating table: %v", err)
		panic(err)
	}
	log.Println("Database migrated successfully.")
}
