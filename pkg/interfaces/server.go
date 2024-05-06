package interfaces

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mikaijun/aquagent/pkg/infrastructure"
	"github.com/mikaijun/aquagent/pkg/infrastructure/repositoryimpl"
	"github.com/mikaijun/aquagent/pkg/interfaces/handler"
	"github.com/mikaijun/aquagent/pkg/usecase"
)

var r *gin.Engine

func Serve(addr string) {
	userRepoImpl := repositoryimpl.NewUserRepositoryImpl(infrastructure.Conn)
	questionRepoImpl := repositoryimpl.NewQuestionRepositoryImpl(infrastructure.Conn)
	userUseCase := usecase.NewUserUseCase(userRepoImpl)
	questionUseCase := usecase.NewQuestionUseCase(questionRepoImpl)
	userHandler := handler.NewUserHandler(userUseCase)
	questionHandler := handler.NewQuestionHandler(questionUseCase)

	r = gin.Default()

	r.POST("/signup", userHandler.HandleSignup)
	r.POST("/login", userHandler.HandleLogin)
	r.GET("/logout", userHandler.HandleLogout)

	secured := r.Group("/secured").Use(Middleware())

	secured.GET("/user", userHandler.HandleFetchUser)
	secured.GET("/questions", questionHandler.HandleGetAll)
	secured.GET("/question/:id", questionHandler.HandleGet)
	secured.POST("/question", questionHandler.HandleCreate)
	secured.PUT("/question/:id", questionHandler.HandleUpdate)

	log.Println("Server running...")
	if err := r.Run(addr); err != nil {
		log.Fatalf("Listen and serve failed. %+v", err)
	}
}
