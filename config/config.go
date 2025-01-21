package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type AuthConfig struct {
	JwtSecret []byte
}

type Config struct {
	Database DatabaseConfig
	Auth     AuthConfig
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf(".env file not found %w\n", err)
	}

	// Get the database config
	port, err := strconv.Atoi(os.Getenv("MYSQL_PORT"))
	if err != nil {
		return nil, fmt.Errorf("invalid db port %w", err)
	}

	config := &Config{
		Database: DatabaseConfig{
			Host:     os.Getenv("MYSQL_HOST"),
			Port:     port,
			User:     os.Getenv("MYSQL_USER"),
			Password: os.Getenv("MYSQL_PASSWORD"),
			Name:     os.Getenv("MYSQL_DBNAME"),
		},
		Auth: AuthConfig{
			JwtSecret: []byte(os.Getenv("JWT_SECRET")),
		},
	}

	return config, nil
}
