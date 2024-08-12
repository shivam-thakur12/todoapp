package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"context"

	"github.com/go-redis/redis/v8"
)

type Config struct {
	Database   DatabaseConfig  `toml:"database"`
	Migrations MigrationConfig `toml:"migrations"`
	Redis      RedisConfig     `toml:"redis"`
}

type DatabaseConfig struct {
	User     string `toml:"user"`
	Password string `toml:"password"`
	Dbname   string `toml:"dbname"`
	Sslmode  string `toml:"sslmode"`
}
type MigrationConfig struct {
	Path string `toml:"path"`
}

type RedisConfig struct {
	Address  string `toml:"address"`
	Password string `toml:"password"`
	DB       int    `toml:"db"`
	CacheKey string `toml:"cache_key"`
}

var db *sql.DB

// Initialize the database connection
func initDB(config Config) {

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

var ctx = context.Background()

func NewRedisClient(config Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     config.Redis.Address,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})
}
