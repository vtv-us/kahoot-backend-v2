package services

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vtv-us/kahoot-backend/internal/constants"
	"github.com/vtv-us/kahoot-backend/internal/repositories"
	"github.com/vtv-us/kahoot-backend/internal/utils"
	"github.com/vtv-us/kahoot-backend/internal/utils/gmail"
)

type GroupService struct {
	DB           repositories.Store
	EmailService *gmail.SendgridService
	Config       *utils.Config
}

func NewGroupService(db repositories.Store, sendgrid *gmail.SendgridService, c *utils.Config) *GroupService {
	return &GroupService{
		DB:           db,
		EmailService: sendgrid,
		Config:       c,
	}
}

type createGroupRequest struct {
	GroupName string `json:"group_name"`
}

func (s *GroupService) CreateGroup(ctx *gin.Context) {
	userID := ctx.GetString(constants.Token_USER_ID)

	var req createGroupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	arg := repositories.CreateGroupParams{
		GroupID:   uuid.NewString(),
		GroupName: req.GroupName,
		CreatedBy: userID,
	}
	group, err := s.DB.CreateGroup(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	err = s.DB.AddMemberToGroup(ctx, repositories.AddMemberToGroupParams{
		GroupID: group.GroupID,
		UserID:  userID,
		Role:    constants.Role_OWNER,
		Status:  constants.UserGroupStatus_JOINED,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, group)
}

func (s *GroupService) ListGroupCreatedByUser(ctx *gin.Context) {
	userID := ctx.GetString(constants.Token_USER_ID)

	groups, err := s.DB.ListGroupCreatedByUser(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, groups)
}

func (s *GroupService) ListGroupJoinedByUser(ctx *gin.Context) {
	userID := ctx.GetString(constants.Token_USER_ID)

	groups, err := s.DB.ListGroupJoined(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, groups)
}

type showMemberRequest struct {
	GroupID string `json:"group_id"`
}

func (s *GroupService) ShowGroupMember(ctx *gin.Context) {
	req := showMemberRequest{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	members, err := s.DB.ListMemberInGroup(ctx, req.GroupID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, members)
}

type assignRoleRequest struct {
	GroupID string `json:"group_id"`
	UserID  string `json:"user_id"`
	Role    string `json:"role"`
}

func (s *GroupService) AssignRole(ctx *gin.Context) {
	userID := ctx.GetString(constants.Token_USER_ID)

	var req assignRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	s.checkOwnerPermission(ctx, req.GroupID)

	err := s.DB.UpdateMemberRole(ctx, repositories.UpdateMemberRoleParams{
		GroupID: req.GroupID,
		UserID:  req.UserID,
		Role:    req.Role,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	if req.Role == constants.Role_OWNER {
		err = s.DB.UpdateMemberRole(ctx, repositories.UpdateMemberRoleParams{
			GroupID: req.GroupID,
			UserID:  userID,
			Role:    constants.Role_CO_OWNER,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
			return
		}
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse())
}

func (s *GroupService) checkOwnerPermission(ctx *gin.Context, groupID string) {
	userID := ctx.GetString(constants.Token_USER_ID)

	role, err := s.DB.GetRoleInGroup(ctx, repositories.GetRoleInGroupParams{
		GroupID: groupID,
		UserID:  userID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	if role != constants.Role_OWNER && role != constants.Role_CO_OWNER {
		ctx.JSON(http.StatusForbidden, utils.ErrorResponse(fmt.Errorf("you don't have permission to do this action")))
		return
	}
}

type joinGroupRequest struct {
	GroupID string `uri:"groupid"`
}

func (s *GroupService) JoinGroup(ctx *gin.Context) {
	userID := ctx.GetString(constants.Token_USER_ID)

	var req joinGroupRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	err := s.DB.AddMemberToGroup(ctx, repositories.AddMemberToGroupParams{
		GroupID: req.GroupID,
		UserID:  userID,
		Role:    constants.Role_MEMBER,
		Status:  constants.UserGroupStatus_JOINED,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse())
}
