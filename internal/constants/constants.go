package constants

const (
	Token_USER_ID = "user_id"
	Token_EMAIL   = "email"

	Role_OWNER        = "owner"
	Role_CO_OWNER     = "co-owner"
	Role_MEMBER       = "member"
	Role_COLLABORATOR = "collaborator"

	UserGroupStatus_JOINED   = "joined"
	UserGroupStatus_PENDING  = "pending"
	UserGroupStatus_DECLINED = "declined"

	Cookies_ACCESS_TOKEN = "cookieAccess"

	SocketParticipantStatus_ACTIVE = "active"
	SocketParticipantStatus_LEFT   = "left"

	QuestionType_MULTIPLE_CHOICE = "multiple-choice"
	QuestionType_PARAGRAPH       = "paragraph"
	QuestionType_HEADING         = "heading"
)
