package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mikaijun/anli/pkg/util"
)

type MyJWTClaims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

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

		claims := &MyJWTClaims{}
		token, err := jwt.ParseWithClaims(signedToken, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return util.GetJWTSecret(), nil
		})

		if err != nil {
			v, _ := err.(*jwt.ValidationError)
			switch v.Errors {
			case jwt.ValidationErrorSignatureInvalid:
				err = errors.New("signature validation failed")
			case jwt.ValidationErrorExpired:
				err = errors.New("token is expired")
			default:
				err = errors.New("token is invalid")
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		c.SetCookie("userId", claims.ID, 60*60*24, "/", "localhost", false, true)
		c.Next()
	}
}
