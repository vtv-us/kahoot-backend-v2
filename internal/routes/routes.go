package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/google"
	"github.com/vtv-us/kahoot-backend/internal/services"
	"github.com/vtv-us/kahoot-backend/internal/utils"
)

func InitRoutes(server *services.Server) *gin.Engine {
	route := gin.Default()
	a := services.InitAuthMiddleware(server.AuthService)
	route.Use(a.CORSMiddleware)

	route.LoadHTMLGlob("template/*.html")

	route.POST("/auth/register", server.AuthService.Register)
	route.POST("/auth/login", server.AuthService.Login)
	route.GET("/auth/verify/:email/:code", server.AuthService.Verify)
	route.POST("/auth/resend/:email", server.AuthService.ResendEmail)

	route.GET("/auth/:provider", server.AuthService.LoginProvider)
	route.GET("/auth/:provider/callback", server.AuthService.ProviderCallback)

	auth := route.Group("/auth")
	auth.Use(a.AuthRequired)
	auth.POST("/change-password", server.AuthService.ChangePassword)
	auth.GET("/refresh", server.AuthService.Refresh)

	group := route.Group("/group")
	group.Use(a.AuthRequired)
	group.POST("/", server.GroupService.CreateGroup)
	group.GET("/", server.GroupService.ListGroupCreatedByUser)
	group.GET("/link/:groupid", server.GroupService.GetGroupLink)
	group.GET("/joined", server.GroupService.ListGroupJoinedByUser)
	group.GET("/:groupid/member", server.GroupService.ShowGroupMember)
	group.POST("/role", server.GroupService.AssignRole)
	group.POST("/:groupid", server.GroupService.JoinGroup)
	group.POST("/:groupid/leave", server.GroupService.LeaveGroup)
	group.POST("/kick", server.GroupService.KickMember)
	group.POST("/invite", server.GroupService.InviteMember)

	user := route.Group("/user")
	user.Use(a.AuthRequired)
	user.GET("/profile", server.UserService.GetProfile)
	user.GET("/profile/:userid", server.UserService.GetProfileByUserID)
	user.POST("/profile", server.UserService.UpdateProfile)
	user.POST("/avatar", server.UserService.UploadAvatar)

	return route
}

func InitGoth(config *utils.Config) {
	goth.UseProviders(
		facebook.New(config.FBKey, config.FBSecret, config.FrontendAddress+"/auth/facebook/callback"),
		google.New(config.GGKey, config.GGSecret, config.FrontendAddress+"/auth/google/callback"),
	)
}
