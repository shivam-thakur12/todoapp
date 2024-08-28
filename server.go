package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/contribsys/faktory/client"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"context"

	"github.com/go-redis/redis/v8"
)

var db *sql.DB
var faktoryClient *client.Client

// Initialize the database connection
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

var ctx = context.Background()

func NewRedisClient(config Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     config.Redis.Address,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})
}

func initFaktory(config Config) *client.Client {
	os.Setenv("FAKTORY_URL", config.Faktory.URL)
	var err error
	faktoryClient, err = client.Open()
	if err != nil {
		log.Fatalf("Error connecting to Faktory: %v", err)
	}
	return faktoryClient
}
