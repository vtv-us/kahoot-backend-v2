package services

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vtv-us/kahoot-backend/internal/constants"
	"github.com/vtv-us/kahoot-backend/internal/repositories"
	"github.com/vtv-us/kahoot-backend/internal/utils"
	"github.com/vtv-us/kahoot-backend/internal/utils/cloudinary"
)

type UserService struct {
	DB                repositories.Store
	CloudinaryService *cloudinary.CloudinaryService
	Config            *utils.Config
}

func NewUserService(db repositories.Store, cloudinary *cloudinary.CloudinaryService, c *utils.Config) *UserService {
	return &UserService{
		DB:                db,
		CloudinaryService: cloudinary,
		Config:            c,
	}
}

type getProfileResponse struct {
	User userResponse `json:"user"`
}

func (s *UserService) GetProfile(ctx *gin.Context) {
	userID := ctx.GetString(constants.Token_USER_ID)

	user, err := s.DB.GetUser(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, getProfileResponse{
		User: userResponse{
			UserID:     user.UserID,
			Name:       user.Name,
			Email:      user.Email,
			AvatarUrl:  user.AvatarUrl,
			Verified:   user.Verified,
			GoogleID:   user.GoogleID.String,
			FacebookID: user.FacebookID.String,
		},
	})
}

type uploadAvatarResponse struct {
	ImageUrl string       `json:"image_url"`
	User     userResponse `json:"user"`
}

func (s *UserService) UploadAvatar(ctx *gin.Context) {
	userID := ctx.GetString(constants.Token_USER_ID)

	formFile, _, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(fmt.Errorf("invalid file")))
		return
	}
	defer formFile.Close()

	uploadUrl, err := s.CloudinaryService.UploadImage(formFile)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(fmt.Errorf("failed to upload image: %w", err)))
		return
	}

	user, err := s.DB.UpdateAvatarUrl(ctx, repositories.UpdateAvatarUrlParams{
		AvatarUrl: uploadUrl,
		UserID:    userID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(fmt.Errorf("failed to update avatar url: %w", err)))
		return
	}

	ctx.JSON(http.StatusOK, uploadAvatarResponse{
		ImageUrl: uploadUrl,
		User: userResponse{
			UserID:     user.UserID,
			Email:      user.Email,
			Name:       user.Name,
			AvatarUrl:  uploadUrl,
			Verified:   user.Verified,
			GoogleID:   user.GoogleID.String,
			FacebookID: user.FacebookID.String,
		},
	})
}

type updateProfileRequest struct {
	Name string `json:"name" binding:"required"`
}

type updateProfileResponse struct {
	User userResponse `json:"user"`
}

func (s *UserService) UpdateProfile(ctx *gin.Context) {
	userID := ctx.GetString(constants.Token_USER_ID)

	var req updateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	user, err := s.DB.UpdateProfile(ctx, repositories.UpdateProfileParams{
		UserID: userID,
		Name:   req.Name,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(fmt.Errorf("failed to update profile: %w", err)))
		return
	}

	ctx.JSON(http.StatusOK, updateProfileResponse{
		User: userResponse{
			UserID:     user.UserID,
			Email:      user.Email,
			Name:       user.Name,
			AvatarUrl:  user.AvatarUrl,
			Verified:   user.Verified,
			GoogleID:   user.GoogleID.String,
			FacebookID: user.FacebookID.String,
		},
	})
}
