package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/vtv-us/kahoot-backend/internal/entities"
	"github.com/vtv-us/kahoot-backend/internal/repositories"
)

func (s *SlideService) SaveChatMsg(slideID, username, content string) error {
	_, err := s.DB.SaveChat(context.Background(), repositories.SaveChatParams{
		ID:       uuid.NewString(),
		SlideID:  slideID,
		Username: username,
		Content:  content,
	})
	return err
}

func (s *SlideService) GetChatMsgs(slideID string) ([]entities.ChatMsg, error) {
	chatMsg, err := s.DB.GetChatBySlide(context.Background(), slideID)
	if err != nil {
		return nil, err
	}

	chatMsgEnt := make([]entities.ChatMsg, 0, len(chatMsg))
	for _, msg := range chatMsg {
		chatMsgEnt = append(chatMsgEnt, msg.ChatMsg)
	}
	return chatMsgEnt, nil
}
