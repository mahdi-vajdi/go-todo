package auth

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
	"todo/configs"
)

type AuthService interface {
	RegisterUser(credentials Credentials) error
	Login(credentials Credentials) (*LoginResponse, error)
	ValidateToken(tokenString string) (int64, error)
}

type Service struct {
	repository Repository
	cfg        *config.AuthConfig
}

func NewService(repository Repository, config *config.AuthConfig) AuthService {
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

func (s *Service) ValidateToken(tokenString string) (int64, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method for the jwt
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unauthorized")
		}
		return s.cfg.JwtSecret, nil
	})
	if err != nil || !token.Valid {
		return 0, errors.New("unauthorized")
	}

	userIdString, err := token.Claims.GetSubject()
	if err != nil {
		return 0, fmt.Errorf("error extracting user ID from token: %w", err)
	}

	userId, _ := strconv.ParseInt(userIdString, 10, 64)

	return userId, nil
}

func (s *Service) generateToken(userId int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "Go-Todo",
		Subject:   strconv.FormatInt(userId, 10),
		Audience:  nil,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		NotBefore: jwt.NewNumericDate(time.Now()),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	tokenString, err := token.SignedString(s.cfg.JwtSecret)
	if err != nil {
		return "", fmt.Errorf("error signing token: %w", err)
	}

	return tokenString, nil
}
