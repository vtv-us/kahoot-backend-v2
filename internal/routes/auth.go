package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vtv-us/kahoot-backend/internal/services"
)

func InitRoutes(server *services.Server) *gin.Engine {
	route := gin.Default()

	route.POST("/auth/register", server.AuthService.Register)

	return route
}
