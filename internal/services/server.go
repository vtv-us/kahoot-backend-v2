package services

import (
	"github.com/vtv-us/kahoot-backend/internal/repositories"
	"github.com/vtv-us/kahoot-backend/internal/utils"
)

type Server struct {
	AuthService  *AuthService
	GroupService *GroupService
}

func NewServer(store repositories.Store, c *utils.Config) *Server {

	jwt := utils.JwtWrapper{
		SecretKey: c.JwtSecretKey,
		Issuer:    "go-grpc-auth-svc",
	}

	authService := NewAuthService(store, &jwt, c)
	groupService := NewGroupService(store, c)

	return &Server{
		AuthService:  authService,
		GroupService: groupService,
	}
}
