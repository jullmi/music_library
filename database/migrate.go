package database

import (
	"fmt"

	_ "github.com/lib/pq"  
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"

)

func RunMigrations(databaseUrl string, migrationPath string) error {
	mPath := migrationPath

	m, err := migrate.New(mPath, databaseUrl)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %v", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to create migrate instance: %v", err)
	}

	fmt.Println("Migrations applied successfully!")
	return nil
}
