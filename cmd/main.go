package main

import (
	"fmt"
	"log"
	"music_library/config"
	"music_library/db"
)

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

	migrationsPath := "db/migrations"
	err := db.RunMigrations(databaseURL, migrationsPath)

	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Println("Application started successfully!")

}
