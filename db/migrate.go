package db

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations (databaseUrl string, migrationPath string) error {
	mPath := fmt.Sprintf("file://%s", migrationPath)

	m, err := migrate.New(mPath, databaseUrl)
	if err != nil {
		return fmt.Errorf("Failed to create migrate instance: %v", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("Failed to create migrate instance: %v", err)
	}

	fmt.Println("Migrations applied successfully!")
	return nil
}