package main

import (
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Database   DatabaseConfig  `toml:"database"`
	Migrations MigrationConfig `toml:"migrations"`
	Faktory    FaktoryConfig   `toml:"faktory"`
}

type DatabaseConfig struct {
	User     string `toml:"user"`
	Password string `toml:"password"`
	Dbname   string `toml:"dbname"`
	Sslmode  string `toml:"sslmode"`
	Host     string `toml:"host"`
}
type MigrationConfig struct {
	Path string `toml:"path"`
}

type FaktoryConfig struct {
	URL      string `toml:"url"`
	Password string `toml:"password"`
}

func initDBConfig(config Config) string {

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s host=%s",
		config.Database.User, config.Database.Password, config.Database.Dbname, config.Database.Sslmode, config.Database.Host)
	return connStr

}

func initMigrationConfig(config Config) (string, string) {
	// Adjust connection string format for migrations
	migrationConnStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		config.Database.User, config.Database.Password, config.Database.Host, config.Database.Dbname, config.Database.Sslmode)
	return migrationConnStr, config.Migrations.Path
}

func initConfig() Config {
	var config Config
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		log.Fatal(err)
	}
	return config
}
