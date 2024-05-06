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
	HandleGetAll(c *gin.Context)
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
			ID        int64  `json:"id"`
			UserID    int64  `json:"user_id"`
			Title     string `json:"title"`
			Content   string `json:"content"`
			FilePath  string `json:"file_path"`
			CreatedAt string `json:"created_at"`
			UpdatedAt string `json:"updated_at"`
			DeletedAt string `json:"deleted_at"`
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
		FilePath:  "",
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
		DeletedAt: "",
	}

	question, err = h.useCase.Create(c.Request.Context(), question)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, &response{
		ID:      question.ID,
		Title:   question.Title,
		Content: question.Content,
	})
}

func (h *questionHandler) HandleGetAll(c *gin.Context) {
	userId, err := util.FindUserIdByCookie(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	questions, err := h.useCase.GetAll(c.Request.Context(), userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, questions)
}
