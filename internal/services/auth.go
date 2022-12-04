package services

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/lib/pq"
	"github.com/oov/gothic"
	"github.com/vtv-us/kahoot-backend/internal/constants"
	"github.com/vtv-us/kahoot-backend/internal/repositories"
	"github.com/vtv-us/kahoot-backend/internal/utils"
	"github.com/vtv-us/kahoot-backend/internal/utils/gmail"
)

type AuthService struct {
	DB           repositories.Store
	EmailService *gmail.SendgridService
	JWT          *utils.JwtWrapper
	Config       *utils.Config
}

func NewAuthService(db repositories.Store, sendgrid *gmail.SendgridService, jwt *utils.JwtWrapper, c *utils.Config) *AuthService {
	return &AuthService{
		DB:           db,
		EmailService: sendgrid,
		JWT:          jwt,
		Config:       c,
	}
}

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

type registerResponse struct {
	UserID     string `json:"user_id"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	Verified   bool   `json:"verified"`
	GoogleID   string `json:"google_id"`
	FacebookID string `json:"facebook_id"`
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
		UserID:       uuid.NewString(),
		Email:        req.Email,
		Name:         req.Name,
		Password:     hashedPassword,
		Verified:     false,
		VerifiedCode: uuid.NewString(),
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

	err = s.EmailService.SendEmailForVerified(user.Email, user.VerifiedCode)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(fmt.Errorf("user create but send email failed, pls send again: %w", err)))
		return
	}

	rsp := registerResponse{
		UserID:     user.UserID,
		Email:      user.Email,
		Name:       user.Name,
		Verified:   user.Verified,
		GoogleID:   user.GoogleID.String,
		FacebookID: user.FacebookID.String,
	}

	ctx.JSON(http.StatusOK, rsp)
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type userResponse struct {
	UserID     string `json:"user_id"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	AvatarUrl  string `json:"avatar_url"`
	Verified   bool   `json:"verified"`
	GoogleID   string `json:"google_id"`
	FacebookID string `json:"facebook_id"`
}

type loginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         userResponse `json:"user"`
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

	if !user.Verified {
		ctx.JSON(http.StatusForbidden, utils.ErrorResponse(fmt.Errorf("user not verified")))
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
		User: userResponse{
			UserID:     user.UserID,
			Email:      user.Email,
			Name:       user.Name,
			AvatarUrl:  user.AvatarUrl,
			Verified:   user.Verified,
			GoogleID:   user.GoogleID.String,
			FacebookID: user.FacebookID.String,
		},
	}

	ctx.JSON(http.StatusOK, rsp)
}

type refreshResponse struct {
	AccessToken string `json:"access_token"`
}

func (s *AuthService) Refresh(ctx *gin.Context) {
	email := ctx.GetString(constants.Token_EMAIL)

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

func (s *AuthService) LoginProvider(ctx *gin.Context) {
	err := gothic.BeginAuth(ctx.Param("provider"), ctx.Writer, ctx.Request)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

func (s *AuthService) ProviderCallback(ctx *gin.Context) {
	gUser, err := gothic.CompleteAuth(ctx.Param("provider"), ctx.Writer, ctx.Request)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// check if user already exist
	exist := true
	user, err := s.DB.GetUserByEmail(ctx, gUser.Email)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			exist = false
		} else {
			ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
			return
		}
	}

	verifyCode := uuid.NewString()
	// if not, create new user
	if !exist {
		arg := repositories.CreateUserParams{
			UserID:       uuid.NewString(),
			Email:        gUser.Email,
			Name:         gUser.Name,
			Password:     "",
			Verified:     true,
			VerifiedCode: verifyCode,
		}
		provider := ctx.Param("provider")
		if provider == "google" {
			arg.GoogleID = utils.NullString(gUser.UserID)
		}
		if provider == "facebook" {
			arg.FacebookID = utils.NullString(gUser.UserID)
		}
		user, err = s.DB.CreateUser(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
			return
		}

		// if yes, update user
	} else {
		arg := repositories.UpdateSocialIDParams{
			Email: user.Email,
		}
		provider := ctx.Param("provider")
		if provider == "google" {
			arg.GoogleID = utils.NullString(gUser.UserID)
		}
		if provider == "facebook" {
			arg.FacebookID = utils.NullString(gUser.UserID)
		}
		user, err = s.DB.UpdateSocialID(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
			return
		}
		user, err = s.DB.UpdateVerifiedCode(ctx, repositories.UpdateVerifiedCodeParams{
			UserID:       user.UserID,
			VerifiedCode: verifyCode,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
			return
		}
	}

	ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/auth/callback/%s/%s", s.Config.FrontendAddress, user.UserID, verifyCode))
}

type loginCallbackRequest struct {
	UserID string `uri:"user_id"`
	Code   string `uri:"code"`
}

func (s *AuthService) LoginCallback(ctx *gin.Context) {
	var req loginCallbackRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	user, err := s.DB.GetUser(ctx, req.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	if user.VerifiedCode != req.Code {
		ctx.JSON(http.StatusForbidden, utils.ErrorResponse(fmt.Errorf("invalid code")))
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
		User: userResponse{
			UserID:     user.UserID,
			Email:      user.Email,
			Name:       user.Name,
			AvatarUrl:  user.AvatarUrl,
			Verified:   user.Verified,
			GoogleID:   user.GoogleID.String,
			FacebookID: user.FacebookID.String,
		},
	}
	ctx.JSON(http.StatusOK, rsp)
}

type verifyRequest struct {
	Email string `uri:"email" binding:"required,email"`
	Code  string `uri:"code" binding:"required"`
}

func (s *AuthService) Verify(ctx *gin.Context) {
	var req verifyRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
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

	if user.Verified {
		ctx.JSON(http.StatusForbidden, utils.ErrorResponse(fmt.Errorf("user already verified")))
		return
	}

	if (user.VerifiedCode) != req.Code {
		ctx.JSON(http.StatusForbidden, utils.ErrorResponse(fmt.Errorf("invalid code")))
		return
	}

	user, err = s.DB.Verify(ctx, req.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.HTML(http.StatusOK, "success.html", gin.H{
		"content": s.Config.FrontendAddress,
	})
}

type resendEmail struct {
	Email string `uri:"email" binding:"required"`
}

func (s *AuthService) ResendEmail(ctx *gin.Context) {
	var req resendEmail
	if err := ctx.ShouldBindUri(&req); err != nil {
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

	if user.Verified {
		ctx.JSON(http.StatusForbidden, utils.ErrorResponse(fmt.Errorf("user already verified")))
		return
	}

	// send email
	err = s.EmailService.SendEmailForVerified(req.Email, user.VerifiedCode)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "email sent",
	})
}

type changePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

func (s *AuthService) ChangePassword(ctx *gin.Context) {
	userID := ctx.GetString(constants.Token_USER_ID)
	var req changePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	user, err := s.DB.GetUser(ctx, userID)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusForbidden, utils.ErrorResponse(fmt.Errorf("user not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	if !user.Verified {
		ctx.JSON(http.StatusForbidden, utils.ErrorResponse(fmt.Errorf("user not verified")))
		return
	}

	if user.Password != "" {
		if err := utils.CheckPassword(req.OldPassword, user.Password); err != nil {
			ctx.JSON(http.StatusForbidden, utils.ErrorResponse(fmt.Errorf("wrong password")))
			return
		}
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	user, err = s.DB.UpdatePassword(ctx, repositories.UpdatePasswordParams{
		UserID:   userID,
		Password: hashedPassword,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse())
}
