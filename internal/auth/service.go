package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
	"todo/config"
)

type Service struct {
	db  *sql.DB
	cfg *config.AuthConfig
}

type LoginResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

func NewService(db *sql.DB, config *config.AuthConfig) *Service {
	return &Service{db: db, cfg: config}
}

func (s *Service) RegisterUser(credentials Credentials) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	query := `INSERT INTO users (email, password, created_at) VALUES (?, ?, ?)`
	_, err = s.db.Exec(query, credentials.Email, hashedPassword, time.Now())
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	return nil
}

func (s *Service) Login(credentials Credentials) (*LoginResponse, error) {
	var user User
	var hashedPassword string

	query := `SELECT id, email, password, created_at FROM users WHERE email = ?`
	err := s.db.QueryRow(query, credentials.Email).Scan(&user.Id, &user.Email, &hashedPassword, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("invalid credentials")
		}
		return nil, fmt.Errorf("error getting user from database: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(credentials.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := s.generateToken(user.Id)
	if err != nil {
		return nil, fmt.Errorf("error generating token: %w", err)
	}

	return &LoginResponse{
		User:  user,
		Token: token,
	}, nil
}

func (s *Service) generateToken(userId int64) (string, error) {
	claims := jwt.MapClaims{"user_id": userId, "exp": time.Now().Add(24 * time.Hour).Unix(), "iat": time.Now().Unix()}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(s.cfg.JwtSecret)
	if err != nil {
		return "", fmt.Errorf("error signing token: %w", err)
	}

	return tokenString, nil
}
