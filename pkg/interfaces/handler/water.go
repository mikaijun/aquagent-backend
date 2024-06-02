package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mikaijun/aquagent/pkg/domain/model"
	"github.com/mikaijun/aquagent/pkg/usecase"
	"github.com/mikaijun/aquagent/pkg/util"
)

type WaterHandler interface {
	HandleSearch(c *gin.Context)
	HandleCreate(c *gin.Context)
	HandleDelete(c *gin.Context)
}

type waterHandler struct {
	useCase usecase.WaterUseCase
}

func NewWaterHandler(waterUseCase usecase.WaterUseCase) WaterHandler {
	return &waterHandler{
		useCase: waterUseCase,
	}
}

func (h *waterHandler) HandleSearch(c *gin.Context) {
	userId, err := util.FindUserIdByCookie(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	filters := make(map[string]interface{})

	date := c.Query("date")
	if date != "" {
		filters["date"] = date
	}

	month := c.Query("month")
	if month != "" {
		filters["month"] = month
	}

	waters, err := h.useCase.Search(c.Request.Context(), userId, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, waters)
}

func (h *waterHandler) HandleCreate(c *gin.Context) {
	type (
		request struct {
			Volume int64 `json:"volume" binding:"required"`
		}
		response struct {
			ID        int64  `json:"id"`
			Volume    int64  `json:"volume" binding:"required"`
			CreatedAt string `json:"created_at"`
			UpdatedAt string `json:"updated_at"`
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

	water := &model.Water{
		Volume:    requestBody.Volume,
		UserID:    userId,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	water, err = h.useCase.Create(c.Request.Context(), water)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, &response{
		ID:        water.ID,
		Volume:    water.Volume,
		CreatedAt: water.CreatedAt,
		UpdatedAt: water.UpdatedAt,
	})
}

func (h *waterHandler) HandleDelete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = util.FindUserIdByCookie(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.useCase.Delete(c.Request.Context(), id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "water delete successful"})
}
