package services

import (
	"fmt"

	"github.com/vtv-us/kahoot-backend/internal/constants"
	"github.com/vtv-us/kahoot-backend/internal/repositories"

	"github.com/gin-gonic/gin"
)

func checkGroupPermission(ctx *gin.Context, db repositories.Store, groupID string, opt string) error {
	userID := ctx.GetString(constants.Token_USER_ID)

	role, err := db.GetRoleInGroup(ctx, repositories.GetRoleInGroupParams{
		GroupID: groupID,
		UserID:  userID,
	})
	if err != nil {
		return err
	}

	if role != constants.Role_OWNER && role != constants.Role_CO_OWNER {
		return fmt.Errorf("you don't have permission to do this action")
	}

	if opt == constants.Role_OWNER {
		if role != constants.Role_OWNER {
			return fmt.Errorf("you don't have permission to do this action")
		}
	}

	return nil
}

func checkQuestionPermission(ctx *gin.Context, db repositories.Store, questionID string) error {
	userID := ctx.GetString(constants.Token_USER_ID)

	isAllowed, err := db.CheckQuestionPermission(ctx, repositories.CheckQuestionPermissionParams{
		ID:    questionID,
		Owner: userID,
	})
	if err != nil {
		return err
	}

	if !isAllowed {
		return fmt.Errorf("you are not allowed to access this question")
	}

	return nil
}

func checkSlidePermission(ctx *gin.Context, db repositories.Store, slideID string) error {
	userID := ctx.GetString(constants.Token_USER_ID)

	isAllowed, err := db.CheckSlidePermission(ctx, repositories.CheckSlidePermissionParams{
		ID:    slideID,
		Owner: userID,
	})
	if err != nil {
		return err
	}

	if !isAllowed {
		return fmt.Errorf("you are not allowed to access this slide")
	}

	return nil
}

func checkAnswerPermission(ctx *gin.Context, db repositories.Store, answerID string) error {
	userID := ctx.GetString(constants.Token_USER_ID)

	isAllowed, err := db.CheckAnswerPermission(ctx, repositories.CheckAnswerPermissionParams{
		ID:    answerID,
		Owner: userID,
	})
	if err != nil {
		return err
	}

	if !isAllowed {
		return fmt.Errorf("you are not allowed to access this answer")
	}

	return nil
}
