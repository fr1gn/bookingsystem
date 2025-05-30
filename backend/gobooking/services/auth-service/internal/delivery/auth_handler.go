package delivery

import (
	"context"
	"github.com/fr1gn/bookingsystem/backend/gobooking/services/auth-service/auth"
	"github.com/fr1gn/bookingsystem/backend/gobooking/services/auth-service/internal/service"
)

type AuthHandler struct {
	auth.UnimplementedAuthServiceServer
	Service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{Service: service}
}

func (h *AuthHandler) RegisterUser(ctx context.Context, req *auth.RegisterRequest) (*auth.AuthResponse, error) {
	user, err := h.Service.RegisterUser(ctx, req.FullName, req.Email, req.Password)
	if err != nil {
		log.Println("RegisterUser error:", err)
		return &auth.AuthResponse{Message: err.Error()}, nil
	}
	log.Println("User registered:", user.Email)
	return &auth.AuthResponse{
		Message: "User registered successfully. Check your email for verification.",
	}, nil
}

func (h *AuthHandler) LoginUser(ctx context.Context, req *auth.LoginRequest) (*auth.AuthResponse, error) {
	access, refresh, err := h.Service.LoginUser(ctx, req.Email, req.Password)
	if err != nil {
		log.Println("LoginUser error:", err)
		return &auth.AuthResponse{Message: err.Error()}, nil
	}
	log.Println("User logged in:", req.Email)
	return &auth.AuthResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		Message:      "Login successful",
	}, nil
}

func (h *AuthHandler) ValidateToken(ctx context.Context, req *auth.TokenRequest) (*auth.ValidateResponse, error) {
	valid, userID := h.Service.ValidateToken(req.Token)

	return &auth.ValidateResponse{
		IsValid: valid,
		UserId:  userID,
	}, nil
}

func (h *AuthHandler) SendVerificationEmail(ctx context.Context, req *auth.EmailRequest) (*auth.GenericResponse, error) {
	err := h.Service.SendVerificationEmail(req.Email)
	if err != nil {
		return &auth.GenericResponse{Message: err.Error()}, nil
	}
	return &auth.GenericResponse{Message: "Verification email sent"}, nil
}
