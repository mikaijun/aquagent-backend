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
	HandleGet(c *gin.Context)
	HandleGetAll(c *gin.Context)
	HandleCreate(c *gin.Context)
	HandleUpdate(c *gin.Context)
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

func (h *waterHandler) HandleGet(c *gin.Context) {
	type response struct {
		ID        int64  `json:"id"`
		Volume    int64  `json:"volume"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	water, err := h.useCase.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userId, err := util.FindUserIdByCookie(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if water.UserID != userId {
		c.JSON(http.StatusForbidden, gin.H{"error": "you are not allowed to access this data"})
		return
	}

	c.JSON(http.StatusOK, &response{
		ID:        water.ID,
		Volume:    water.Volume,
		CreatedAt: water.CreatedAt,
		UpdatedAt: water.UpdatedAt,
	})
}

func (h *waterHandler) HandleGetAll(c *gin.Context) {
	userId, err := util.FindUserIdByCookie(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	waters, err := h.useCase.GetAll(c.Request.Context(), userId)
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

func (h *waterHandler) HandleUpdate(c *gin.Context) {
	type (
		request struct {
			Volume int64 `json:"volume" binding:"required"`
		}
		response struct {
			ID        int64  `json:"id"`
			Volume    int64  `json:"volume" binding:"required"`
			UpdatedAt string `json:"updated_at"`
		}
	)

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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
		ID:        id,
		Volume:    requestBody.Volume,
		UserID:    userId,
		UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	water, err = h.useCase.Update(c.Request.Context(), water)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, &response{
		ID:        water.ID,
		Volume:    water.Volume,
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
