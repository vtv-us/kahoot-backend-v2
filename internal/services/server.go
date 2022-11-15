package services

import (
	"github.com/vtv-us/kahoot-backend/internal/repositories"
	"github.com/vtv-us/kahoot-backend/internal/utils"
)

type Server struct {
	AuthService *AuthService
}

func NewServer(store repositories.Store, c *utils.Config) *Server {

	jwt := utils.JwtWrapper{
		SecretKey: c.JwtSecretKey,
		Issuer:    "go-grpc-auth-svc",
	}

	authService := NewAuthService(store, &jwt, c)

	return &Server{
		AuthService: authService,
	}
}
