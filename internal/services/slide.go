package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vtv-us/kahoot-backend/internal/constants"
	"github.com/vtv-us/kahoot-backend/internal/repositories"
	"github.com/vtv-us/kahoot-backend/internal/utils"
)

type SlideService struct {
	DB     repositories.Store
	Config *utils.Config
}

func NewSlideService(db repositories.Store, c *utils.Config) *SlideService {
	return &SlideService{
		DB:     db,
		Config: c,
	}
}

type createSlideRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func (s *SlideService) CreateSlide(ctx *gin.Context) {
	userID := ctx.GetString(constants.Token_USER_ID)

	var req createSlideRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	slide, err := s.DB.CreateSlide(ctx, repositories.CreateSlideParams{
		ID:      uuid.NewString(),
		Owner:   userID,
		Title:   req.Title,
		Content: req.Content,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, slide)
}

func (s *SlideService) GetSlideByUserID(ctx *gin.Context) {
	userID := ctx.GetString(constants.Token_USER_ID)

	slides, err := s.DB.GetSlidesByOwner(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, slides)
}

type getSlideByIDRequest struct {
	SlideID string `uri:"slide_id" binding:"required"`
}

func (s *SlideService) GetSlideByID(ctx *gin.Context) {
	var req getSlideByIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	slide, err := s.DB.GetSlide(ctx, req.SlideID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, slide)
}

type updateSlideRequest struct {
	SlideID string `json:"slide_id" binding:"required"`
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func (s *SlideService) UpdateSlide(ctx *gin.Context) {
	var req updateSlideRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	if err := checkSlidePermission(ctx, s.DB, req.SlideID); err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse(err))
		return
	}

	slide, err := s.DB.UpdateSlide(ctx, repositories.UpdateSlideParams{
		ID:      req.SlideID,
		Title:   req.Title,
		Content: req.Content,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, slide)
}

type deleteSlideRequest struct {
	SlideID string `uri:"slide_id" binding:"required"`
}

func (s *SlideService) DeleteSlide(ctx *gin.Context) {
	var req deleteSlideRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	err := checkSlidePermission(ctx, s.DB, req.SlideID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse(err))
		return
	}

	err = s.DB.DeleteSlideTx(ctx, req.SlideID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse())
}

type addCollaboratorRequest struct {
	SlideID string `json:"slide_id" binding:"required"`
	UserID  string `json:"user_id" binding:"required"`
}

func (s *SlideService) AddCollaborator(ctx *gin.Context) {
	var req addCollaboratorRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	if err := checkSlidePermission(ctx, s.DB, req.SlideID); err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse(err))
		return
	}

	err := s.DB.AddCollab(ctx, repositories.AddCollabParams{
		SlideID: req.SlideID,
		UserID:  req.UserID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse())
}

type getCollaboratorBySlideIDRequest struct {
	SlideID string `uri:"slide_id" binding:"required"`
}

func (s *SlideService) GetCollaboratorBySlideID(ctx *gin.Context) {
	var req getCollaboratorBySlideIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	collabs, err := s.DB.ListCollabBySlide(ctx, req.SlideID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, collabs)
}

type getCollaboratorByUserIDRequest struct {
	UserID string `uri:"user_id" binding:"required"`
}

func (s *SlideService) GetCollaboratorByUserID(ctx *gin.Context) {
	var req getCollaboratorByUserIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	collabs, err := s.DB.ListCollab(ctx, req.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, collabs)
}

type removeCollaboratorRequest struct {
	SlideID string `json:"slide_id" binding:"required"`
	UserID  string `json:"user_id" binding:"required"`
}

func (s *SlideService) RemoveCollaborator(ctx *gin.Context) {
	var req removeCollaboratorRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	err := checkSlidePermission(ctx, s.DB, req.SlideID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse(err))
		return
	}

	err = s.DB.RemoveCollab(ctx, repositories.RemoveCollabParams{
		SlideID: req.SlideID,
		UserID:  req.UserID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse())
}
