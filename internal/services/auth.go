package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/vtv-us/kahoot-backend/internal/repositories"
	db "github.com/vtv-us/kahoot-backend/internal/repositories"
	"github.com/vtv-us/kahoot-backend/internal/utils"
)

type AuthService struct {
	DB  repositories.Store
	JWT *utils.JwtWrapper
}

func NewAuthService(db repositories.Store, jwt *utils.JwtWrapper) *AuthService {
	return &AuthService{
		DB:  db,
		JWT: jwt,
	}
}

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

type registerResponse struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Verified bool   `json:"verified"`
	GoogleID string `json:"google_id"`
}

func (s *AuthService) Register(ctx *gin.Context) {
	var req registerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	arg := db.CreateUserParams{
		UserID:   uuid.NewString(),
		Email:    req.Email,
		Name:     req.Name,
		Password: hashedPassword,
	}

	user, err := s.DB.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation", "foreign_key_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	rsp := registerResponse{
		UserID:   user.UserID,
		Email:    user.Email,
		Name:     user.Name,
		Verified: user.Verified,
		GoogleID: user.GoogleID.String,
	}

	ctx.JSON(http.StatusOK, rsp)
}

// type loginRequest struct {
// 	Email    string `json:"email" binding:"required,email"`
// 	Password string `json:"password" binding:"required,min=6"`
// }

// type loginResponse struct {
// 	AccessToken  string `json:"access_token"`
// 	RefreshToken string `json:"refresh_token"`
// }

// func (s *AuthService) login(ctx *gin.Context) {
// 	var req loginRequest
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}

// 	user, err := s.DB.GetUserByEmail(ctx, req.Email)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	if err := utils.CheckPassword(req.Password, user.Password); err != nil {
// 		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("incorrect password")))
// 		return
// 	}

// 	accessToken, err := utils.GenerateAccessToken(user.UserID)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	refreshToken, err := utils.GenerateRefreshToken(user.UserID)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	rsp := loginResponse{
// 		AccessToken:  accessToken,
// 		RefreshToken: refreshToken,
// 	}

// 	ctx.JSON(http.StatusOK, rsp)
// }
