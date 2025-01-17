package database

import (
	"database/sql"
	"fmt"
	"time"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func NewConnection(config Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", config.User, config.Password, config.Host, config.Port, config.DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening databse: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the datbase: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}
