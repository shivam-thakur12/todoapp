package server

import (
	"TODO/pkg/config"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/contribsys/faktory/client"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/go-redis/redis/v8"
)

var DB *sql.DB
var faktoryClient *client.Client

// Initialize the database connection
func InitDB(configg config.Config) {

	connStr := config.InitDBConfig(configg)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database connected successfully!")

	runMigrations(configg)
}

func runMigrations(configg config.Config) {

	migrationConnStr, migrationPath := config.InitMigrationConfig(configg)

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

func NewRedisClient(configg config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     configg.Redis.Address,
		Password: configg.Redis.Password,
		DB:       configg.Redis.DB,
	})
}

func InitFaktory(configg config.Config) *client.Client {
	os.Setenv("FAKTORY_URL", configg.Faktory.URL)
	var err error
	faktoryClient, err = client.Open()
	if err != nil {
		log.Fatalf("Error connecting to Faktory: %v", err)
	}
	return faktoryClient
}
