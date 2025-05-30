package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/fr1gn/bookingsystem/backend/gobooking/services/auth-service/internal/model"
	"github.com/fr1gn/bookingsystem/backend/gobooking/services/auth-service/internal/repository"
	"github.com/fr1gn/bookingsystem/backend/gobooking/services/auth-service/pkg/email"
	"github.com/fr1gn/bookingsystem/backend/gobooking/services/auth-service/shared"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var JWTSecret = []byte("elbobritobandito228")

type AuthService struct {
	UserRepo *repository.UserRepo
}

func NewAuthService(repo *repository.UserRepo) *AuthService {
	return &AuthService{UserRepo: repo}
}

func (s *AuthService) SendVerificationEmail(emailAddr string) error {
	code := generateCode()
	return email.SendVerificationEmail(emailAddr, code)
}

func (s *AuthService) RegisterUser(ctx context.Context, fullName, emailAddr, password string) (*model.User, error) {
	log.Println("RegisterUser called for:", emailAddr)
	existing, _ := s.UserRepo.FindByEmail(ctx, emailAddr)
	if existing != nil {
		log.Println("User already exists:", emailAddr)
		return nil, errors.New("user already exists")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Password hashing failed:", err)
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
		log.Println("CreateUser failed:", err)
		return nil, err
	}
	go email.SendVerificationEmail(user.Email, generateCode())
	log.Println("User created:", emailAddr)
	return user, nil
}

func (s *AuthService) LoginUser(ctx context.Context, emailAddr, password string) (string, string, error) {
	log.Println("LoginUser called for:", emailAddr)
	userJSON, err := shared.GetCache("user:" + emailAddr)
	var user *model.User
	if err == nil {
		user = &model.User{}
		if err := json.Unmarshal([]byte(userJSON), user); err != nil {
			log.Println("Cache unmarshal error:", err)
			return "", "", err
		}
	} else {
		user, err = s.UserRepo.FindByEmail(ctx, emailAddr)
		if err != nil {
			log.Println("User not found:", emailAddr)
			return "", "", errors.New("user not found")
		}
		data, _ := json.Marshal(user)
		_ = shared.SetCache("user:"+emailAddr, string(data), 10*time.Minute)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		log.Println("Invalid password for user:", emailAddr)
		return "", "", errors.New("invalid credentials")
	}
	accessToken, err := generateJWT(user.ID, 15*time.Minute)
	if err != nil {
		log.Println("Access token generation failed:", err)
		return "", "", err
	}
	refreshToken, err := generateJWT(user.ID, 7*24*time.Hour)
	if err != nil {
		log.Println("Refresh token generation failed:", err)
		return "", "", err
	}
	log.Println("User logged in successfully:", emailAddr)
	return accessToken, refreshToken, nil
}

// ----- TOKEN VALIDATION -----
// Валидация токена
func (s *AuthService) ValidateToken(token string) (bool, string) {
	claims, err := parseJWT(token)
	if err != nil {
		log.Println("Error parsing token:", err)
		return false, ""
	}
	return true, claims.UserID
}

// ----- JWT LOGIC -----
// Генерация JWT
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func generateJWT(userID string, duration time.Duration) (string, error) {
	key := JWTSecret
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
	key := JWTSecret
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

// ----- SIMPLE EMAIL VERIFICATION CODE -----
// Генерация кода подтверждения
func generateCode() string {
	return time.Now().Format("150405")
}
