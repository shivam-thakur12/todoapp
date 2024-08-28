package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var db *sql.DB

func initDB(config Config) {
	connStr := initDBConfig(config)
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database connected successfully!")
	runMigrations(config)
}
func runMigrations(config Config) {

	migrationConnStr, migrationPath := initMigrationConfig(config)

	// Run migrations
	// Connection string to the database
	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationPath), // Source path to migration files from config
		migrationConnStr)
	if err != nil {
		log.Fatal(err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	fmt.Println("Database migrated successfully!")
}
