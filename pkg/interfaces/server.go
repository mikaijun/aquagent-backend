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
	waterRepoImpl := repositoryimpl.NewWaterRepositoryImpl(infrastructure.Conn)
	userUseCase := usecase.NewUserUseCase(userRepoImpl)
	waterUseCase := usecase.NewWaterUseCase(waterRepoImpl)
	userHandler := handler.NewUserHandler(userUseCase)
	waterHandler := handler.NewWaterHandler(waterUseCase)

	r = gin.Default()

	r.POST("/signup", userHandler.HandleSignup)
	r.POST("/login", userHandler.HandleLogin)
	r.GET("/logout", userHandler.HandleLogout)

	secured := r.Group("/secured").Use(Middleware())

	secured.GET("/user", userHandler.HandleFetchUser)
	secured.GET("/waters", waterHandler.HandleGetAll)
	secured.GET("/water/:id", waterHandler.HandleGet)
	secured.POST("/water", waterHandler.HandleCreate)
	secured.PUT("/water/:id", waterHandler.HandleUpdate)

	log.Println("Server running...")
	if err := r.Run(addr); err != nil {
		log.Fatalf("Listen and serve failed. %+v", err)
	}
}
