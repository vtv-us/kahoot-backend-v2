package repositories

import (
	"github.com/vtv-us/kahoot-backend/internal/entities"
)

type User struct {
	entities.User
}

type Group struct {
	entities.Group
}

type UserGroup struct {
	entities.UserGroup
}

type Slide struct {
	entities.Slide
}

type Question struct {
	entities.Question
}

type Answer struct {
	entities.Answer
}

type AnswerHistory struct {
	entities.AnswerHistory
}

type UserQuestion struct {
	entities.UserQuestion
}

type ChatMsg struct {
	entities.ChatMsg
}
