package routes

import (
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/google"
	"github.com/vtv-us/kahoot-backend/internal/services"
	"github.com/vtv-us/kahoot-backend/internal/utils"
)

func InitRoutes(server *services.Server, socket *socketio.Server) *gin.Engine {
	// gin.SetMode(gin.ReleaseMode)
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
	route.GET("/auth/callback/:user_id/:code", server.AuthService.LoginCallback)

	auth := route.Group("/auth")
	auth.Use(a.AuthRequired)
	auth.POST("/change-password", server.AuthService.ChangePassword)
	auth.GET("/refresh", server.AuthService.Refresh)

	group := route.Group("/group")
	group.Use(a.AuthRequired)
	group.POST("", server.GroupService.CreateGroup)
	group.GET("", server.GroupService.ListGroupCreatedByUser)
	group.GET("/:id", server.GroupService.GetGroupByID)
	group.GET("/link/:groupid", server.GroupService.GetGroupLink)
	group.GET("/joined", server.GroupService.ListGroupJoinedByUser)
	group.GET("/member/:groupid", server.GroupService.ShowGroupMember)
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

	slide := route.Group("/slide")
	slide.GET("/:slide_id", server.SlideService.GetSlideByID)
	slide.Use(a.AuthRequired)
	slide.POST("", server.SlideService.CreateSlide)
	slide.GET("", server.SlideService.GetSlideByUserID)
	slide.PUT("", server.SlideService.UpdateSlide)
	slide.DELETE("/:slide_id", server.SlideService.DeleteSlide)

	question := route.Group("/question")
	question.GET("/:question_id", server.QuestionService.GetQuestionByID)
	question.GET("/slide/:slide_id", server.QuestionService.GetQuestionBySlideID)
	question.Use(a.AuthRequired)
	question.POST("", server.QuestionService.CreateQuestion)
	question.PUT("", server.QuestionService.UpdateQuestion)
	question.DELETE("/:question_id", server.QuestionService.DeleteQuestion)

	answer := route.Group("/answer")
	answer.GET("/:answer_id", server.AnswerService.GetAnswerByID)
	answer.GET("/question/:question_id", server.AnswerService.GetAnswerByQuestionID)
	answer.Use(a.AuthRequired)
	answer.POST("", server.AnswerService.CreateAnswer)
	answer.PUT("", server.AnswerService.UpdateAnswer)
	answer.DELETE("/:answer_id", server.AnswerService.DeleteAnswer)

	// route.GET("/socket.io/*any", gin.WrapH(socket))
	// route.POST("/socket.io/*any", gin.WrapH(socket))
	// route.GET("/test", func(c *gin.Context) {
	// 	http.ServeFile(c.Writer, c.Request, "index.html")
	// })

	return route
}

func InitGoth(config *utils.Config) {
	goth.UseProviders(
		facebook.New(config.FBKey, config.FBSecret, config.ServerAddress+"/auth/facebook/callback"),
		google.New(config.GGKey, config.GGSecret, config.ServerAddress+"/auth/google/callback"),
	)
}
