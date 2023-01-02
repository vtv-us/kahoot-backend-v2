package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vtv-us/kahoot-backend/internal/repositories"
	"github.com/vtv-us/kahoot-backend/internal/utils"
)

type QuestionService struct {
	DB     repositories.Store
	Config *utils.Config
}

func NewQuestionService(db repositories.Store, c *utils.Config) *QuestionService {
	return &QuestionService{
		DB:     db,
		Config: c,
	}
}

type createQuestionRequest struct {
	SlideID         string `json:"slide_id" binding:"required"`
	Index           int16  `json:"index" binding:"required"`
	RawQuestion     string `json:"raw_question" binding:"required"`
	Meta            string `json:"meta"`
	LongDescription string `json:"long_description"`
	Type            string `json:"type" binding:"required" oneof:"multiple-choice paragraph heading"`
}

func (s *QuestionService) CreateQuestion(ctx *gin.Context) {
	var req createQuestionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	if err := checkSlidePermission(ctx, s.DB, req.SlideID); err != nil {
		ctx.JSON(http.StatusForbidden, utils.ErrorResponse(err))
		return
	}

	question, err := s.DB.CreateQuestion(ctx, repositories.CreateQuestionParams{
		ID:              uuid.NewString(),
		SlideID:         req.SlideID,
		Index:           req.Index,
		RawQuestion:     req.RawQuestion,
		Meta:            req.Meta,
		LongDescription: req.LongDescription,
		Type:            req.Type,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, question)
}

type getQuestionBySlideIDRequest struct {
	SlideID string `uri:"slide_id" binding:"required"`
}

func (s *QuestionService) GetQuestionBySlideID(ctx *gin.Context) {
	var req getQuestionBySlideIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	question, err := s.DB.GetQuestionsBySlide(ctx, req.SlideID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, question)
}

type getQuestionByIDRequest struct {
	QuestionID string `uri:"question_id" binding:"required"`
}

func (s *QuestionService) GetQuestionByID(ctx *gin.Context) {
	var req getQuestionByIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	question, err := s.DB.GetQuestion(ctx, req.QuestionID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, question)
}

type updateQuestionRequest struct {
	QuestionID      string `json:"question_id" binding:"required"`
	Index           int16  `json:"index" binding:"required"`
	RawQuestion     string `json:"raw_question" binding:"required"`
	Meta            string `json:"meta"`
	LongDescription string `json:"long_description"`
	Type            string `json:"type" binding:"required" oneof:"multiple-choice paragraph heading"`
}

func (s *QuestionService) UpdateQuestion(ctx *gin.Context) {
	var req updateQuestionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	err := checkQuestionPermission(ctx, s.DB, req.QuestionID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse(err))
		return
	}

	question, err := s.DB.UpdateQuestion(ctx, repositories.UpdateQuestionParams{
		ID:              req.QuestionID,
		Index:           req.Index,
		RawQuestion:     req.RawQuestion,
		Meta:            req.Meta,
		LongDescription: req.LongDescription,
		Type:            req.Type,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, question)
}

type deleteQuestionRequest struct {
	QuestionID string `uri:"question_id" binding:"required"`
}

func (s *QuestionService) DeleteQuestion(ctx *gin.Context) {
	var req deleteQuestionRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	err := checkQuestionPermission(ctx, s.DB, req.QuestionID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse(err))
		return
	}

	err = s.DB.DeleteQuestionTx(ctx, req.QuestionID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse())
}
