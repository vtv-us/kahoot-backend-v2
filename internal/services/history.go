package services

import (
	"context"

	"github.com/vtv-us/kahoot-backend/internal/entities"
	"github.com/vtv-us/kahoot-backend/internal/repositories"
)

func (s *SlideService) SaveAnswerHistory(username, slideID, questionID, answerID string) error {
	_, err := s.DB.UpsertAnswerHistory(context.Background(), repositories.UpsertAnswerHistoryParams{
		Username:   username,
		SlideID:    slideID,
		QuestionID: questionID,
		AnswerID:   answerID,
	})
	return err
}

func (s *SlideService) ListAnswerHistoryBySlideID(slideID string) ([]entities.AnswerHistory, error) {
	res, err := s.DB.ListAnswerHistoryBySlideID(context.Background(), slideID)
	if err != nil {
		return nil, err
	}

	answers := make([]entities.AnswerHistory, 0, len(res))
	for _, answer := range res {
		answers = append(answers, answer.AnswerHistory)
	}
	return answers, nil
}

func (s *SlideService) ListAnswerHistoryByQuestionID(questionID string) ([]entities.AnswerHistory, error) {
	res, err := s.DB.ListAnswerHistoryByQuestionID(context.Background(), questionID)
	if err != nil {
		return nil, err
	}

	answers := make([]entities.AnswerHistory, 0, len(res))
	for _, answer := range res {
		answers = append(answers, answer.AnswerHistory)
	}
	return answers, nil
}

func (s *SlideService) ListAnswerHistoryByAnswerID(answerID string) ([]entities.AnswerHistory, error) {
	res, err := s.DB.ListAnswerHistoryByAnswerID(context.Background(), answerID)
	if err != nil {
		return nil, err
	}

	answers := make([]entities.AnswerHistory, 0, len(res))
	for _, answer := range res {
		answers = append(answers, answer.AnswerHistory)
	}
	return answers, nil
}

type AnswerCount struct {
	AnswerID string
	Count    int
}

func (s *SlideService) CountAnswerByQuestionID(questionID string) ([]AnswerCount, error) {
	res, err := s.DB.CountAnswerByQuestionID(context.Background(), questionID)
	if err != nil {
		return nil, err
	}

	answers := make([]AnswerCount, 0, len(res))
	for _, answer := range res {
		answers = append(answers, AnswerCount{
			AnswerID: answer.AnswerID,
			Count:    int(answer.Count),
		})
	}

	return answers, nil
}
