package server

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mikaijun/anli/pkg/infrastructure"
	"github.com/mikaijun/anli/pkg/infrastructure/repositoryimpl"
	"github.com/mikaijun/anli/pkg/interfaces/api/handler"
	"github.com/mikaijun/anli/pkg/interfaces/api/middleware"
	"github.com/mikaijun/anli/pkg/usecase"
)

var r *gin.Engine

func Serve(addr string) {
	userRepoImpl := repositoryimpl.NewUserRepositoryImpl(infrastructure.Conn)
	userUseCase := usecase.NewUserUseCase(userRepoImpl)
	userHandler := handler.NewHandler(userUseCase)

	r = gin.Default()

	r.POST("/signup", userHandler.HandleSignup)
	r.POST("/login", userHandler.HandleLogin)
	r.GET("/logout", userHandler.HandleLogout)

	secured := r.Group("/secured").Use(middleware.Auth())
	secured.GET("/user", userHandler.HandleFetchUser)

	log.Println("Server running...")
	if err := r.Run(addr); err != nil {
		log.Fatalf("Listen and serve failed. %+v", err)
	}
}
