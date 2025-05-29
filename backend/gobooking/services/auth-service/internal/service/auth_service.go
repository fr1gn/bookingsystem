package service

import (
	"auth-service/internal/model"
	"auth-service/internal/repository"
	"auth-service/pkg/email"
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo *repository.UserRepo
}

func NewAuthService(repo *repository.UserRepo) *AuthService {
	return &AuthService{UserRepo: repo}
}

// ----- REGISTER -----

func (s *AuthService) RegisterUser(ctx context.Context, fullName, emailAddr, password string) (*model.User, error) {
	existing, _ := s.UserRepo.FindByEmail(ctx, emailAddr)
	if existing != nil {
		return nil, errors.New("user already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		FullName:      fullName,
		Email:         emailAddr,
		PasswordHash:  string(hash),
		EmailVerified: false,
	}

	err = s.UserRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	// Send verification code
	code := generateCode()
	go email.SendVerificationEmail(user.Email, code) // без блокировки

	return user, nil
}

// ----- LOGIN -----

func (s *AuthService) LoginUser(ctx context.Context, emailAddr, password string) (string, string, error) {
	user, err := s.UserRepo.FindByEmail(ctx, emailAddr)
	if err != nil {
		return "", "", errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	accessToken, err := generateJWT(user.ID, 15*time.Minute)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := generateJWT(user.ID, 7*24*time.Hour)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// ----- TOKEN VALIDATION -----

func (s *AuthService) ValidateToken(token string) (bool, string) {
	claims, err := parseJWT(token)
	if err != nil {
		return false, ""
	}
	return true, claims.UserID
}

// ----- JWT LOGIC -----

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func generateJWT(userID string, duration time.Duration) (string, error) {
	key := []byte(os.Getenv("JWT_SECRET"))
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(key)
}

func parseJWT(tokenStr string) (*Claims, error) {
	key := []byte(os.Getenv("JWT_SECRET"))
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

// ----- SIMPLE EMAIL VERIFICATION CODE -----

func generateCode() string {
	return time.Now().Format("150405") // Пример: "153207"
}
