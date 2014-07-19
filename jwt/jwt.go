package jwt_middleware

import (
	"errors"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func JWT(secret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := jwt.ParseFromRequest(c.Request, func(token *jwt.Token) ([]byte, error) {
			return secret, nil
		})

		if err == nil && token.Valid {
			c.Next()
		} else {
			c.Fail(401, errors.New("Unauthorized"))
		}
	}
}
