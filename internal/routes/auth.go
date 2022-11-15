package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vtv-us/kahoot-backend/internal/services"
)

func InitRoutes(server *services.Server) *gin.Engine {
	route := gin.Default()
	a := services.InitAuthMiddleware(server.AuthService)

	route.POST("/auth/register", server.AuthService.Register)
	route.POST("/auth/login", server.AuthService.Login)

	auth := route.Group("/auth")
	auth.Use(a.AuthRequired)
	auth.GET("/refresh", server.AuthService.Refresh)

	return route
}
