package services

import (
	"database/sql"
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
	GroupName   string `json:"group_name" binding:"required"`
	Description string `json:"description"`
}

func (s *GroupService) CreateGroup(ctx *gin.Context) {
	userID := ctx.GetString(constants.Token_USER_ID)

	var req createGroupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	arg := repositories.CreateGroupParams{
		GroupID:     uuid.NewString(),
		GroupName:   req.GroupName,
		CreatedBy:   userID,
		Description: req.Description,
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

	groups, err := s.DB.ListGroupOwned(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, groups)
}

type getGroupByIDRequest struct {
	GroupID string `uri:"id" binding:"required"`
}

func (s *GroupService) GetGroupByID(ctx *gin.Context) {
	var req getGroupByIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	group, err := s.DB.GetGroup(ctx, req.GroupID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, group)
}

type getGroupLinkRequest struct {
	GroupID string `uri:"groupid" binding:"required"`
}

type getGroupLinkResponse struct {
	GroupLink string `json:"group_link"`
}

func (s *GroupService) GetGroupLink(ctx *gin.Context) {
	userID := ctx.GetString(constants.Token_USER_ID)
	var req getGroupLinkRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	err := s.checkUserInGroup(ctx, req.GroupID, userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	link := utils.GenLink(s.Config.FrontendAddress, req.GroupID, userID)

	ctx.JSON(http.StatusOK, getGroupLinkResponse{
		GroupLink: link,
	})
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
	GroupID string `uri:"groupid" binding:"required"`
}

func (s *GroupService) ShowGroupMember(ctx *gin.Context) {
	req := showMemberRequest{}

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	members, err := s.DB.ListMemberInGroup(ctx, req.GroupID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	// sort by role
	// owner -> co-owner -> member
	owner := make([]repositories.ListMemberInGroupRow, 0)
	coOwner := make([]repositories.ListMemberInGroupRow, 0)
	member := make([]repositories.ListMemberInGroupRow, 0)
	for _, m := range members {
		if m.Role == constants.Role_OWNER {
			owner = append(owner, m)
		} else if m.Role == constants.Role_CO_OWNER {
			coOwner = append(coOwner, m)
		} else {
			member = append(member, m)
		}
	}

	ctx.JSON(http.StatusOK, append(append(owner, coOwner...), member...))
}

type assignRoleRequest struct {
	GroupID string `json:"group_id" binding:"required"`
	UserID  string `json:"user_id" binding:"required"`
	Role    string `json:"role" binding:"required,oneof=owner co-owner member"`
}

func (s *GroupService) AssignRole(ctx *gin.Context) {
	userID := ctx.GetString(constants.Token_USER_ID)

	var req assignRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	err := s.checkOwnerPermission(ctx, req.GroupID, constants.Role_OWNER)
	if err != nil {
		ctx.JSON(http.StatusForbidden, utils.ErrorResponse(err))
		return
	}

	err = s.DB.UpdateMemberRole(ctx, repositories.UpdateMemberRoleParams{
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

func (s *GroupService) checkOwnerPermission(ctx *gin.Context, groupID string, opt string) error {
	userID := ctx.GetString(constants.Token_USER_ID)

	role, err := s.DB.GetRoleInGroup(ctx, repositories.GetRoleInGroupParams{
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

type joinGroupRequest struct {
	GroupID string `uri:"groupid" binding:"required"`
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

type leaveGroupRequest struct {
	GroupID string `uri:"groupid" binding:"required"`
}

func (s *GroupService) LeaveGroup(ctx *gin.Context) {
	userID := ctx.GetString(constants.Token_USER_ID)

	var req leaveGroupRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	userGroup, err := s.DB.GetUserGroup(ctx, repositories.GetUserGroupParams{
		GroupID: req.GroupID,
		UserID:  userID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	if userGroup.Role == constants.Role_OWNER {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(fmt.Errorf("you can't leave group if you are owner")))
		return
	}

	err = s.DB.RemoveMemberFromGroup(ctx, repositories.RemoveMemberFromGroupParams{
		GroupID: req.GroupID,
		UserID:  userID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse())
}

type kickMemberRequest struct {
	GroupID string `json:"group_id" binding:"required"`
	UserID  string `json:"user_id" binding:"required"`
}

func (s *GroupService) KickMember(ctx *gin.Context) {
	userID := ctx.GetString(constants.Token_USER_ID)

	var req kickMemberRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	if userID == req.UserID {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(fmt.Errorf("you can't kick yourself")))
		return
	}

	err := s.checkOwnerPermission(ctx, req.GroupID, "")
	if err != nil {
		ctx.JSON(http.StatusForbidden, utils.ErrorResponse(err))
		return
	}

	err = s.DB.RemoveMemberFromGroup(ctx, repositories.RemoveMemberFromGroupParams{
		GroupID: req.GroupID,
		UserID:  req.UserID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse())
}

type inviteMemberRequest struct {
	GroupID string `json:"group_id" binding:"required"`
	Email   string `json:"email" binding:"required,email"`
}

func (s *GroupService) InviteMember(ctx *gin.Context) {
	inviterID := ctx.GetString(constants.Token_USER_ID)
	var req inviteMemberRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	err := s.checkUserInGroup(ctx, req.GroupID, inviterID)
	if err != nil {
		ctx.JSON(http.StatusForbidden, utils.ErrorResponse(err))
		return
	}

	emails, err := s.DB.ListEmailInGroup(ctx, req.GroupID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	for _, email := range emails {
		if email == req.Email {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(fmt.Errorf("this email is already in group")))
			return
		}
	}

	inviter, err := s.DB.GetUser(ctx, inviterID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	group, err := s.DB.GetGroup(ctx, req.GroupID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	user, err := s.DB.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			err = s.EmailService.SendEmailForInvite(req.Email, req.GroupID, group.GroupName, inviter.Name, inviter.UserID)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(fmt.Errorf("can't send email")))
				return
			}
			ctx.JSON(http.StatusOK, utils.ResponseWithMessage("this email is not registered, we will send an invitation email to this email"))
			return
		}
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	err = s.DB.AddMemberToGroup(ctx, repositories.AddMemberToGroupParams{
		GroupID: req.GroupID,
		UserID:  user.UserID,
		Role:    constants.Role_MEMBER,
		Status:  constants.UserGroupStatus_PENDING,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	err = s.EmailService.SendEmailForInvite(user.Email, req.GroupID, group.GroupName, inviter.Name, inviter.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(fmt.Errorf("can't send email")))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse())
}

func (s *GroupService) checkUserInGroup(ctx *gin.Context, groupID string, userID string) error {
	userIsInGroup, err := s.DB.CheckUserInGroup(ctx, repositories.CheckUserInGroupParams{
		GroupID: groupID,
		UserID:  userID,
	})
	if err != nil {
		return err
	}
	if !userIsInGroup {
		return fmt.Errorf("you are not in this group")
	}

	return nil
}

type deleteGroupRequest struct {
	GroupID string `uri:"groupid" binding:"required"`
}

func (s *GroupService) DeleteGroup(ctx *gin.Context) {
	var req deleteGroupRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	err := s.checkOwnerPermission(ctx, req.GroupID, constants.Role_OWNER)
	if err != nil {
		ctx.JSON(http.StatusForbidden, utils.ErrorResponse(err))
		return
	}

	err = s.DB.DeleteGroup(ctx, req.GroupID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse())
}
