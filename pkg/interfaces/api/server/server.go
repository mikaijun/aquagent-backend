package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikaijun/anli/pkg/infrastructure"
	"github.com/mikaijun/anli/pkg/infrastructure/repositoryimpl"
	"github.com/mikaijun/anli/pkg/interfaces/api/handler"
	"github.com/mikaijun/anli/pkg/interfaces/api/middleware"
	"github.com/mikaijun/anli/pkg/usecase"
)

var r *gin.Engine

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

// サーバ起動処理
func Serve(addr string) {
	// 依存性の注入
	userRepoImpl := repositoryimpl.NewRepositoryImpl(infrastructure.Conn)
	userUseCase := usecase.NewUseCase(userRepoImpl)
	userHandler := handler.NewHandler(userUseCase)

	r = gin.Default()

	r.POST("/signup", userHandler.HandleSignup)
	r.POST("/login", userHandler.HandleLogin)
	r.GET("/logout", userHandler.HandleLogout)

	secured := r.Group("/secured").Use(middleware.Auth())
	secured.GET("/ping", Ping)

	log.Println("Server running...")
	if err := r.Run(addr); err != nil {
		log.Fatalf("Listen and serve failed. %+v", err)
	}
}
