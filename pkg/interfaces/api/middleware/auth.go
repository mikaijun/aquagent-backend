package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikaijun/anli/pkg/util"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		signedToken, err := c.Cookie("jwt")

		if signedToken == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("no token set in cookie").Error()})
			c.Abort()
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errors.New("cookie is not found").Error()})
			c.Abort()
			return
		}

		err = util.ValidateToken(signedToken)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Next()
	}
}
