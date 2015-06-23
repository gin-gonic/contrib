package jwt

import (
	"fmt"

	jwt_lib "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Auth(secret string, alg string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := jwt_lib.ParseFromRequest(c.Request, func(token *jwt_lib.Token) (interface{}, error) {

			switch alg {
			case "HS256", "HS384", "HS512":
				if _, ok := token.Method.(*jwt_lib.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
			case "RSA256", "RSA384", "RSA512":
				if _, ok := token.Method.(*jwt_lib.SigningMethodRSA); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
			default:
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			b := ([]byte(secret))
			return b, nil
		})

		if err != nil {
			c.AbortWithError(401, err)
		}
		if !token.Valid {
			c.AbortWithError(401, fmt.Errorf("Invalid Token"))
		}
	}
}
