package services

import (
	"github.com/gin-gonic/gin"
	"github.com/vtv-us/kahoot-backend/internal/repositories"
	"github.com/vtv-us/kahoot-backend/internal/utils"
)

type Server struct {
	AuthService *AuthService
}

func NewServer(store repositories.Store, c *utils.Config) *Server {

	jwt := utils.JwtWrapper{
		SecretKey: c.JWT_SECRET_KET,
		Issuer:    "go-grpc-auth-svc",
	}

	authService := NewAuthService(store, &jwt)

	return &Server{
		AuthService: authService,
	}
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
