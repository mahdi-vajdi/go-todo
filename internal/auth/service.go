package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

func (s *Service) RegisterUser(credentials Credentials) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %s", err)
	}

	query := `INSERT INTO users (email, password, createdAt) VALUES (?, ?, ?)`
	_, err = s.db.Exec(query, credentials.Email, hashedPassword, time.Now())
	if err != nil {
		return fmt.Errorf("error creating user: %s", err)
	}

	return nil
}

func (s *Service) AuthenticateUser(credentials Credentials) (*User, error) {
	var user User
	var hashedPassword string

	query := `SELECT email, password, createdAt FROM users WHERE email = ?`
	err := s.db.QueryRow(query, credentials.Email).Scan(&user.Id, &user.Email, &hashedPassword, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("invalid credentials")
		}
		return nil, fmt.Errorf("error getting user from databse: %s", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(credentials.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return &user, nil

}
