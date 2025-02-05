package auth

import (
	"database/sql"
	"errors"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(email string, password []byte) error {
	query := `INSERT INTO users (email, password) VALUES (?, ?)`
	_, err := r.db.Exec(query, email, password)
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	return nil
}

func (r *Repository) GetUserByEmail(email string) (*User, error) {
	var user User
	query := `SELECT id, email, password, created_at FROM users WHERE email = ?`
	err := r.db.QueryRow(query, email).Scan(&user.Id, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("error getting user from database: %w", err)
	}

	return &user, nil
}
