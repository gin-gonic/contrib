package jwt

import (
	"bytes"
	"net/http"

	jwt_lib "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := jwt_lib.ParseFromRequest(c.Request, func(token *jwt_lib.Token) (interface{}, error) {
			var b bytes.Buffer
			b.Write([]byte(secret))
			return b, nil
		})

		if err != nil {
			c.Fail(401, err)
		}
		if !token.Valid {
			http.Error(c.Writer, "Invalid Token", http.StatusUnauthorized)
		}
	}
}
