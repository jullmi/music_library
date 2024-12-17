package main

import (
	"fmt"
	"log"
	"music_library/config"
	"music_library/database"
)

// @title           Music Library API
// @version         1.0
// @description     This is an API for managing a music library.
// @host            localhost:8080
// @BasePath        /api
func main() {
	config := config.LoadConfig()
	log.Printf("Configuration loaded: %+v\n", config)

	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.DBUser,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		config.DBName,
	)

	// migrationsPath := "database/migrations"
	err := database.RunMigrations(databaseURL, "postgres://postgres:postgres@localhost:5432/music?sslmode=disable")

	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Println("Application started successfully!")

}
