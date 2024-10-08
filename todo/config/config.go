package config

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Database   DatabaseConfig  `toml:"database"`
	Migrations MigrationConfig `toml:"migrations"`
	Redis      RedisConfig     `toml:"redis"`
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

type RedisConfig struct {
	Address  string `toml:"address"`
	Password string `toml:"password"`
	DB       int    `toml:"db"`
	CacheKey string `toml:"cache_key"`
}
type FaktoryConfig struct {
	URL      string `toml:"url"`
	Password string `toml:"password"`
}

func InitDBConfig(config Config) string {

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s host=%s",
		config.Database.User, config.Database.Password, config.Database.Dbname, config.Database.Sslmode, config.Database.Host)
	return connStr

}

func InitMigrationConfig(config Config) (string, string) {
	// Adjust connection string format for migrations
	migrationConnStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		config.Database.User, config.Database.Password, config.Database.Host, config.Database.Dbname, config.Database.Sslmode)
	return migrationConnStr, config.Migrations.Path
}

func InitConfig() Config {
	var config Config

	// Construct the path to the config.toml file in the root TODO directory
	configPath, err := filepath.Abs("../config.toml")
	if err != nil {
		log.Fatalf("Error finding config file: %v", err)
	}
	// Decode the file into the config struct
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		log.Fatal(err)
	}

	return config
}
