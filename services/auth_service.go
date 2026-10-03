package services

import (
	"errors"
	"time"

	"UTS_PBE/models"
	"UTS_PBE/repositories"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("email atau password salah")
	ErrInvalidJWTSecret   = errors.New("JWT secret belum dikonfigurasi")
)

type AuthService struct {
	UserRepo  *repositories.UserRepository
	JWTSecret string
}

func NewAuthService(
	userRepo *repositories.UserRepository,
	jwtSecret string,
) *AuthService {
	return &AuthService{
		UserRepo:  userRepo,
		JWTSecret: jwtSecret,
	}
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

func (s *AuthService) Login(
	email string,
	password string,
) (*LoginResponse, error) {
	if s.JWTSecret == "" {
		return nil, ErrInvalidJWTSecret
	}

	// Cari pengguna berdasarkan email.
	user, err := s.UserRepo.FindByEmail(email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Periksa password dengan hash bcrypt.
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Buat klaim JWT.
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	// Buat dan tandatangani token.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.JWTSecret))
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token: tokenString,
		User:  user,
	}, nil
}
