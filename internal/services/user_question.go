package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/vtv-us/kahoot-backend/internal/entities"
	"github.com/vtv-us/kahoot-backend/internal/repositories"
	"github.com/vtv-us/kahoot-backend/internal/utils"
)

type UserQuestionService struct {
	DB     repositories.Store
	Config *utils.Config
}

func NewUserQuestionService(db repositories.Store, c *utils.Config) *UserQuestionService {
	return &UserQuestionService{
		DB:     db,
		Config: c,
	}
}

type PostQuestionRequest struct {
	SlideID  string
	Username string
	Content  string
}

func (s *UserQuestionService) PostQuestion(ctx context.Context, req PostQuestionRequest) (entities.UserQuestion, error) {
	question, err := s.DB.UpsertUserQuestion(ctx, repositories.UpsertUserQuestionParams{
		QuestionID: uuid.NewString(),
		SlideID:    req.SlideID,
		Username:   req.Username,
		Content:    req.Content,
	})
	if err != nil {
		return entities.UserQuestion{}, err
	}

	return question.UserQuestion, nil
}

func (s *UserQuestionService) ListQuestionBySlideID(ctx context.Context, slideID string) ([]*entities.UserQuestion, error) {
	questions, err := s.DB.ListUserQuestion(ctx, slideID)
	if err != nil {
		return nil, err
	}

	entQuestions := make([]*entities.UserQuestion, len(questions))
	for i, q := range questions {
		entQuestions[i] = &q.UserQuestion
	}

	return entQuestions, nil
}

func (s *UserQuestionService) UpvoteQuestion(ctx context.Context, questionID string) (entities.UserQuestion, error) {
	question, err := s.DB.UpvoteUserQuestion(ctx, questionID)
	if err != nil {
		return entities.UserQuestion{}, err
	}

	return question.UserQuestion, nil
}

func (s *UserQuestionService) MarkQuestionAsAnswered(ctx context.Context, questionID string) (entities.UserQuestion, error) {
	question, err := s.DB.MarkUserQuestionAnswered(ctx, questionID)
	if err != nil {
		return entities.UserQuestion{}, err
	}

	return question.UserQuestion, nil
}
