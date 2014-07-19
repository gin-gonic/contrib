package jwt

import (
	"errors"
	jwt_lib "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := jwt_lib.ParseFromRequest(c.Request, func(token *jwt_lib.Token) ([]byte, error) {
			return []byte(secret), nil
		})

		if err != nil {
			c.Fail(401, errors.New("Unauthorized token"))
		}
	}
}
