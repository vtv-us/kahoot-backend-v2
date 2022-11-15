package services

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vtv-us/kahoot-backend/internal/utils"
)

type AuthMiddlewareConfig struct {
	auth *AuthService
}

func InitAuthMiddleware(svc *AuthService) AuthMiddlewareConfig {
	return AuthMiddlewareConfig{svc}
}

func (c *AuthMiddlewareConfig) AuthRequired(ctx *gin.Context) {
	authorization := ctx.Request.Header.Get("authorization")

	if authorization == "" {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token := strings.Split(authorization, "Bearer ")

	if len(token) < 2 {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse(fmt.Errorf("invalid authorization header")))
		return
	}

	res, err := c.auth.JWT.ValidateToken(token[1])
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse(err))
		return
	}

	fmt.Println(res.Email)

	ctx.Set("user_id", res.UserID)
	ctx.Set("email", res.Email)

	ctx.Next()
}
