package services

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/lib/pq"
	"github.com/vtv-us/kahoot-backend/internal/repositories"
	"github.com/vtv-us/kahoot-backend/internal/utils"
)

type AuthService struct {
	DB     repositories.Store
	JWT    *utils.JwtWrapper
	Config *utils.Config
}

func NewAuthService(db repositories.Store, jwt *utils.JwtWrapper, c *utils.Config) *AuthService {
	return &AuthService{
		DB:     db,
		JWT:    jwt,
		Config: c,
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
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
	}

	arg := repositories.CreateUserParams{
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
				ctx.JSON(http.StatusForbidden, utils.ErrorResponse(fmt.Errorf("user already exist")))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
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

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *AuthService) Login(ctx *gin.Context) {
	var req loginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	user, err := s.DB.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusForbidden, utils.ErrorResponse(fmt.Errorf("user not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	if err := utils.CheckPassword(req.Password, user.Password); err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse(fmt.Errorf("incorrect password")))
		return
	}

	accessToken, err := s.JWT.GenerateToken(user.User, s.Config.AccessTokenExpiredTime)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	refreshToken, err := s.JWT.GenerateToken(user.User, s.Config.RefreshTokenExpiredTime)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	rsp := loginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	ctx.JSON(http.StatusOK, rsp)
}

type refreshResponse struct {
	AccessToken string `json:"access_token"`
}

func (s *AuthService) Refresh(ctx *gin.Context) {
	email := ctx.GetString("email")

	fmt.Println(email)

	user, err := s.DB.GetUserByEmail(ctx, email)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusForbidden, utils.ErrorResponse(fmt.Errorf("user not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	accessToken, err := s.JWT.GenerateToken(user.User, s.Config.AccessTokenExpiredTime)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	rsp := refreshResponse{
		AccessToken: accessToken,
	}

	ctx.JSON(http.StatusOK, rsp)
}
