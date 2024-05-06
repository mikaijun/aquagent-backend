package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mikaijun/anli/pkg/domain/model"
	"github.com/mikaijun/anli/pkg/usecase"
	"github.com/mikaijun/anli/pkg/util"
)

type QuestionHandler interface {
	HandleCreate(c *gin.Context)
}

type questionHandler struct {
	useCase usecase.QuestionUseCase
}

func NewQuestionHandler(questionUseCase usecase.QuestionUseCase) QuestionHandler {
	return &questionHandler{
		useCase: questionUseCase,
	}
}

func (h *questionHandler) HandleCreate(c *gin.Context) {
	type (
		request struct {
			Title   string `json:"title" binding:"required"`
			Content string `json:"content" binding:"required"`
		}
		response struct {
			ID      int64  `json:"id"`
			Title   string `json:"title"`
			Content string `json:"content"`
		}
	)

	requestBody := new(request)

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, err := util.FindUserIdByCookie(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	question := &model.Question{
		Title:     requestBody.Title,
		Content:   requestBody.Content,
		UserID:    userId,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	question, err = h.useCase.Create(c.Request.Context(), question)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, question)
}
