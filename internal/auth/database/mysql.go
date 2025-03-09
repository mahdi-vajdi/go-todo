package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"todo/internal/auth"
)

type MysqlRepository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) auth.Repository {
	return &MysqlRepository{db: db}
}

func (r *MysqlRepository) CreateUser(email string, password []byte) error {
	query := `INSERT INTO users (email, password) VALUES (:email, :password)`

	_, err := r.db.NamedExec(query, map[string]interface{}{
		"email":    email,
		"password": password,
	})
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	return nil
}

func (r *MysqlRepository) GetUserByEmail(email string) (*auth.User, error) {
	var user auth.User
	query := `SELECT id, email, password, created_at FROM users WHERE email = ?`

	err := r.db.Get(&user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("error getting user from database: %w", err)
	}

	return &user, nil
}
