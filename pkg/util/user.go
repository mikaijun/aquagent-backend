package util

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func FindUserIdByCookie(c *gin.Context) (int64, error) {
	userId, err := c.Cookie("userId")

	if userId == "" {
		return 0, errors.New("no userId set in userId")
	}

	if err != nil {
		return 0, errors.New("userId is not found")
	}

	intUserId, _ := strconv.ParseInt(userId, 10, 64)
	return intUserId, nil
}
