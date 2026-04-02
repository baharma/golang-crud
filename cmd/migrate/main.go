package main

import (
	"crud-test/internal/database"
	"log"
)

func main() {
	log.Println("Running database migration...")

	if err := database.Migrate(); err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("Database migration completed")
}
