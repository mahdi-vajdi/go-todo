package database

import (
	"database/sql"
	"fmt"
	"log"
	"todo/config"
)

var DB *sql.DB

func InitDB() {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		config.DBUser, config.DBPassword, config.DBHost, config.DBName)

	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Could not connect to the datbase: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Could not ping the database %v", err)
	}

	log.Println("Connected to the database successfully")
}
