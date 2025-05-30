package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/fr1gn/bookingsystem/backend/gobooking/services/auth-service/internal/model"
	"github.com/fr1gn/bookingsystem/backend/gobooking/services/auth-service/internal/repository"
	"github.com/fr1gn/bookingsystem/backend/gobooking/services/auth-service/pkg/email"
	"github.com/fr1gn/bookingsystem/backend/gobooking/services/auth-service/shared"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// JWTSecret is the secret key used to sign JWT tokens.
var JWTSecret = []byte("elbobritobandito228")

// AuthService holds dependencies for the authentication service.
type AuthService struct {
	UserRepo *repository.UserRepo
}

// NewAuthService creates a new instance of AuthService.
func NewAuthService(repo *repository.UserRepo) *AuthService {
	return &AuthService{UserRepo: repo}
}

// SendVerificationEmail sends a verification code to the user's email.
func (s *AuthService) SendVerificationEmail(emailAddr string) error {
	code := generateCode()
	return email.SendVerificationEmail(emailAddr, code)
}

// RegisterUser registers a new user if the email is not already taken.
func (s *AuthService) RegisterUser(ctx context.Context, fullName, emailAddr, password string) (*model.User, error) {
	log.Println("RegisterUser called for:", emailAddr)

	// Check if user already exists by email.
	_, err := s.UserRepo.FindByEmail(ctx, emailAddr)
	if err == nil {
		// User found - email already in use.
		log.Println("User already exists:", emailAddr)
		return nil, errors.New("user already exists")
	} else if err != mongo.ErrNoDocuments {
		// Some other error occurred when querying the database.
		log.Println("FindByEmail error:", err)
		return nil, err
	}

	// Hash the user's password before storing.
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

	// Insert the new user into the database.
	err = s.UserRepo.CreateUser(ctx, user)
	if err != nil {
		log.Println("CreateUser failed:", err)
		return nil, err
	}

	// Send verification email asynchronously, log error if sending fails.
	go func() {
		err := email.SendVerificationEmail(user.Email, generateCode())
		if err != nil {
			log.Println("Failed to send verification email:", err)
		}
	}()

	log.Println("User created:", emailAddr)
	return user, nil
}

// LoginUser authenticates a user and returns JWT access and refresh tokens.
func (s *AuthService) LoginUser(ctx context.Context, emailAddr, password string) (string, string, error) {
	log.Println("LoginUser called for:", emailAddr)

	// Try to get user from cache first
	userJSON, err := shared.GetCache("user:" + emailAddr)
	var user *model.User
	if err == nil {
		user = &model.User{}
		if err := json.Unmarshal([]byte(userJSON), user); err != nil {
			log.Println("Cache unmarshal error:", err)
			return "", "", err
		}
	} else {
		// If not in cache, fetch from DB
		user, err = s.UserRepo.FindByEmail(ctx, emailAddr)
		if err != nil {
			log.Println("User not found:", emailAddr)
			return "", "", errors.New("user not found")
		}
		data, _ := json.Marshal(user)
		_ = shared.SetCache("user:"+emailAddr, string(data), 10*time.Minute)
	}

	// Verify password hash matches provided password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		log.Println("Invalid password for user:", emailAddr)
		return "", "", errors.New("invalid credentials")
	}

	// Generate JWT access token valid for 15 minutes
	accessToken, err := generateJWT(user.ID, 15*time.Minute)
	if err != nil {
		log.Println("Access token generation failed:", err)
		return "", "", err
	}

	// Generate JWT refresh token valid for 7 days
	refreshToken, err := generateJWT(user.ID, 7*24*time.Hour)
	if err != nil {
		log.Println("Refresh token generation failed:", err)
		return "", "", err
	}

	log.Println("User logged in successfully:", emailAddr)
	return accessToken, refreshToken, nil
}

// ValidateToken validates a JWT token and returns whether it is valid and the user ID.
func (s *AuthService) ValidateToken(token string) (bool, string) {
	claims, err := parseJWT(token)
	if err != nil {
		log.Println("Error parsing token:", err)
		return false, ""
	}
	return true, claims.UserID
}

// Claims defines JWT claims structure.
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// generateJWT creates a signed JWT token with specified expiration.
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

// parseJWT parses and validates a JWT token string.
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

// generateCode generates a simple verification code (timestamp based).
func generateCode() string {
	return time.Now().Format("150405")
}

func (s *AuthService) VerifyEmail(ctx context.Context, email, code string) error {
	savedCode, err := shared.GetCache("verify:" + email)
	if err != nil {
		return errors.New("verification code expired or not found")
	}
	if savedCode != code {
		return errors.New("invalid verification code")
	}

	err = s.UserRepo.VerifyEmail(ctx, email)
	if err != nil {
		return err
	}

	_ = shared.RedisClient.Del(ctx, "verify:"+email).Err()

	return nil
}
