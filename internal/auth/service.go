package auth

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
	"todo/config"
)

type Service struct {
	repository *Repository
	cfg        *config.AuthConfig
}

func NewService(repository *Repository, config *config.AuthConfig) *Service {
	return &Service{repository: repository, cfg: config}
}

type LoginResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

func (s *Service) RegisterUser(credentials Credentials) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	err = s.repository.CreateUser(credentials.Email, hashedPassword)
	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}

func (s *Service) Login(credentials Credentials) (*LoginResponse, error) {
	user, err := s.repository.GetUserByEmail(credentials.Email)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := s.generateToken(user.Id)
	if err != nil {
		return nil, fmt.Errorf("error generating token: %w", err)
	}

	return &LoginResponse{
		User:  *user,
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
