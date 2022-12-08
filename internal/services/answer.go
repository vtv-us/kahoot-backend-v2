package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vtv-us/kahoot-backend/internal/repositories"
	"github.com/vtv-us/kahoot-backend/internal/utils"
)

type AnswerService struct {
	DB     repositories.Store
	Config *utils.Config
}

func NewAnswerService(db repositories.Store, c *utils.Config) *AnswerService {
	return &AnswerService{
		DB:     db,
		Config: c,
	}
}

type createAnswerRequest struct {
	QuestionID string `json:"question_id" binding:"required"`
	Index      int16  `json:"index" binding:"required"`
	RawAnswer  string `json:"raw_answer" binding:"required"`
}

func (s *AnswerService) CreateAnswer(ctx *gin.Context) {
	var req createAnswerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	question, err := s.DB.CreateAnswer(ctx, repositories.CreateAnswerParams{
		ID:         uuid.NewString(),
		QuestionID: req.QuestionID,
		Index:      req.Index,
		RawAnswer:  req.RawAnswer,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, question)
}

type getAnswerByQuestionIDRequest struct {
	QuestionID string `uri:"question_id" binding:"required"`
}

func (s *AnswerService) GetAnswerByQuestionID(ctx *gin.Context) {
	var req getAnswerByQuestionIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	answer, err := s.DB.GetAnswersByQuestion(ctx, req.QuestionID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, answer)
}

type getAnswerByIDRequest struct {
	AnswerID string `uri:"answer_id" binding:"required"`
}

func (s *AnswerService) GetAnswerByID(ctx *gin.Context) {
	var req getAnswerByIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	question, err := s.DB.GetAnswer(ctx, req.AnswerID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, question)
}

type updateAnswerRequest struct {
	AnswerID  string `json:"answer_id" binding:"required"`
	Index     int16  `json:"index" binding:"required"`
	RawAnswer string `json:"raw_answer" binding:"required"`
}

func (s *AnswerService) UpdateAnswer(ctx *gin.Context) {
	var req updateAnswerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	question, err := s.DB.UpdateAnswer(ctx, repositories.UpdateAnswerParams{
		ID:        req.AnswerID,
		Index:     req.Index,
		RawAnswer: req.RawAnswer,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, question)
}

type deleteAnswerRequest struct {
	AnswerID string `uri:"answer_id" binding:"required"`
}

func (s *AnswerService) DeleteAnswer(ctx *gin.Context) {
	var req deleteAnswerRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	// TODO: check ownership

	err := s.DB.DeleteAnswer(ctx, req.AnswerID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse())
}
