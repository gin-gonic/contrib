package jwt

import (
	jwt_lib "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
)

func Auth(secret interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := request.ParseFromRequest(c.Request, request.OAuth2Extractor, func(token *jwt_lib.Token) (interface{}, error) {
			return secret, nil
		})

		if err != nil {
			c.AbortWithError(401, err)
		}
		if c, ok := token.Claims.(jwt_lib.MapClaims); ok {
			c.Set("JWT", c)
		}
	}
}
