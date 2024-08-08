package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Config struct {
	Database   DatabaseConfig  `toml:"database"`
	Migrations MigrationConfig `toml:"migrations"`
}
type MigrationConfig struct {
	Path string `toml:"path"`
}
type DatabaseConfig struct {
	User     string `toml:"user"`
	Password string `toml:"password"`
	Dbname   string `toml:"dbname"`
	Sslmode  string `toml:"sslmode"`
}

var db *sql.DB

// Initialize the database connection
func initDB() {
	var config Config

	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		log.Fatal(err)
	}

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s",
		config.Database.User, config.Database.Password, config.Database.Dbname, config.Database.Sslmode)

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

	config.runMigrations()

}
func (config Config) runMigrations() {
	// Adjust connection string format for migrations
	migrationConnStr := fmt.Sprintf("postgres://%s:%s@localhost/%s?sslmode=%s",
		config.Database.User, config.Database.Password, config.Database.Dbname, config.Database.Sslmode)

	// Run migrations
	// Connection string to the database
	m, err := migrate.New(
		fmt.Sprintf("file://%s", config.Migrations.Path), // Source path to migration files from config
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
