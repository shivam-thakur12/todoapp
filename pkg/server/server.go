<<<<<<<< HEAD:todo/server.go
package todo

import (
	"TODO/todo/config"
========
package server

import (
	"TODO/pkg/config"
>>>>>>>> a07f97f (worker with refactoring code):pkg/server/server.go
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
<<<<<<<< HEAD:todo/server.go
	DB, _ = sql.Open("postgres", connStr)
	// if err != nil {
	// 	log.Println("Failed to open database:", err)
	// 	DB = nil // Ensure DB is set to nil on error
	// 	return
	// }
========
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
>>>>>>>> a07f97f (worker with refactoring code):pkg/server/server.go

	err = DB.Ping()
	if err != nil {
		log.Println("Failed to ping database:", err)
		DB = nil // Ensure DB is set to nil on error
		return
	}

	fmt.Println("Database connected successfully!")

<<<<<<<< HEAD:todo/server.go
}

func RunMigrations(configg config.Config) {
========
	runMigrations(configg)
}

func runMigrations(configg config.Config) {
>>>>>>>> a07f97f (worker with refactoring code):pkg/server/server.go

	migrationConnStr, migrationPath := config.InitMigrationConfig(configg)

	// Run migrations
	// Connection string to the database
	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationPath), // Source path to migration files from config
		migrationConnStr)
	if err != nil {
		log.Fatalf("Failed to initialize migrations: %v", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to apply migrations: %v", err)
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
