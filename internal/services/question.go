package services

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vtv-us/kahoot-backend/internal/constants"
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
	SlideID       string `json:"slide_id" binding:"required"`
	RawQuestion   string `json:"raw_question" binding:"required"`
	AnswerA       string `json:"answer_a" binding:"required"`
	AnswerB       string `json:"answer_b" binding:"required"`
	AnswerC       string `json:"answer_c" binding:"required"`
	AnswerD       string `json:"answer_d" binding:"required"`
	CorrectAnswer string `json:"correct_answer" binding:"required,oneof=A B C D"`
}

func (s *QuestionService) CreateQuestion(ctx *gin.Context) {
	var req createQuestionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	question, err := s.DB.CreateQuestion(ctx, repositories.CreateQuestionParams{
		ID:            uuid.NewString(),
		SlideID:       req.SlideID,
		RawQuestion:   req.RawQuestion,
		AnswerA:       req.AnswerA,
		AnswerB:       req.AnswerB,
		AnswerC:       req.AnswerC,
		AnswerD:       req.AnswerD,
		CorrectAnswer: req.CorrectAnswer,
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

type updateQuestionRequest struct {
	QuestionID    string `json:"question_id" binding:"required"`
	RawQuestion   string `json:"raw_question" binding:"required"`
	AnswerA       string `json:"answer_a" binding:"required"`
	AnswerB       string `json:"answer_b" binding:"required"`
	AnswerC       string `json:"answer_c" binding:"required"`
	AnswerD       string `json:"answer_d" binding:"required"`
	CorrectAnswer string `json:"correct_answer" binding:"required,oneof=A B C D"`
}

func (s *QuestionService) UpdateQuestion(ctx *gin.Context) {
	var req updateQuestionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	err := s.checkOwnerPermission(ctx, req.QuestionID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse(err))
		return
	}

	question, err := s.DB.UpdateQuestion(ctx, repositories.UpdateQuestionParams{
		ID:            req.QuestionID,
		RawQuestion:   req.RawQuestion,
		AnswerA:       req.AnswerA,
		AnswerB:       req.AnswerB,
		AnswerC:       req.AnswerC,
		AnswerD:       req.AnswerD,
		CorrectAnswer: req.CorrectAnswer,
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

	err := s.checkOwnerPermission(ctx, req.QuestionID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse(err))
		return
	}

	err = s.DB.DeleteQuestion(ctx, req.QuestionID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse())
}

func (s *QuestionService) checkOwnerPermission(ctx *gin.Context, questionID string) error {
	userID := ctx.GetString(constants.Token_USER_ID)

	owner, err := s.DB.GetOwnerOfQuestion(ctx, questionID)
	if err != nil {
		return err
	}

	if owner != userID {
		return fmt.Errorf("you are not the owner of this question")
	}

	return nil
}
