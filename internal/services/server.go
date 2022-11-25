package services

import (
	"github.com/vtv-us/kahoot-backend/internal/repositories"
	"github.com/vtv-us/kahoot-backend/internal/utils"
	"github.com/vtv-us/kahoot-backend/internal/utils/gmail"
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

	emailSvc := gmail.NewMailService(c)
	authService := NewAuthService(store, &emailSvc, &jwt, c)
	groupService := NewGroupService(store, &emailSvc, c)

	return &Server{
		AuthService:  authService,
		GroupService: groupService,
	}
}
