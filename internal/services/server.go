package services

import (
	"github.com/vtv-us/kahoot-backend/internal/repositories"
	"github.com/vtv-us/kahoot-backend/internal/utils"
	"github.com/vtv-us/kahoot-backend/internal/utils/cloudinary"
	"github.com/vtv-us/kahoot-backend/internal/utils/gmail"
)

type Server struct {
	AuthService     *AuthService
	GroupService    *GroupService
	UserService     *UserService
	SlideService    *SlideService
	QuestionService *QuestionService
	AnswerService   *AnswerService
}

func NewServer(store repositories.Store, c *utils.Config) *Server {

	jwt := utils.JwtWrapper{
		SecretKey: c.JwtSecretKey,
		Issuer:    "go-grpc-auth-svc",
	}

	emailSvc := gmail.NewMailService(c)
	cloudinarySvc, err := cloudinary.NewCloudinaryService(c)
	if err != nil {
		panic(err)
	}
	authService := NewAuthService(store, &emailSvc, &jwt, c)
	groupService := NewGroupService(store, &emailSvc, c)
	userService := NewUserService(store, &cloudinarySvc, c)
	slideService := NewSlideService(store, c)
	questionService := NewQuestionService(store, c)
	answerService := NewAnswerService(store, c)

	return &Server{
		AuthService:     authService,
		GroupService:    groupService,
		UserService:     userService,
		SlideService:    slideService,
		QuestionService: questionService,
		AnswerService:   answerService,
	}
}
