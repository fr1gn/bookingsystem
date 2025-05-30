package handler

import (
	authpb "github.com/fr1gn/bookingsystem/backend/gobooking/api-gateway/auth"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func RegisterAuthRoutes(r *gin.Engine, client authpb.AuthServiceClient) {
	r.POST("/auth/register", func(c *gin.Context) {
		var req struct {
			FullName string `json:"full_name"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		res, err := client.RegisterUser(c, &authpb.RegisterRequest{
			FullName: req.FullName,
			Email:    req.Email,
			Password: req.Password,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": res.Message})
	})

	r.POST("/auth/login", func(c *gin.Context) {
		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		res, err := client.LoginUser(c, &authpb.LoginRequest{
			Email:    req.Email,
			Password: req.Password,
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"access_token":  res.AccessToken,
			"refresh_token": res.RefreshToken,
			"message":       res.Message,
		})
	})
	r.GET("/auth/me", func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")

		validResp, err := client.ValidateToken(c, &authpb.TokenRequest{Token: token})
		if err != nil || !validResp.IsValid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		userId := validResp.UserId
		c.JSON(http.StatusOK, gin.H{"user": gin.H{"id": userId}})
	})

}
