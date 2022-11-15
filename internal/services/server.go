package services

import (
	"github.com/gin-gonic/gin"
	"github.com/vtv-us/kahoot-backend/internal/repositories"
)

type Server struct {
	route *gin.Engine

	DB repositories.Store
}

func NewServer(store repositories.Store) *Server {
	server := &Server{
		DB: store,
	}
	route := gin.Default()

	route.POST("/auth/register", server.register)

	server.route = route

	return server
}

func (server *Server) Start(address string) error {
	return server.route.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
