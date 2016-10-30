package jwtauth

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const (
	bearer_ = "Bearer "
)

// Auth(secret,signMethod, contextKey string)
// secret string ....
// signMethod string name for sign method
// contextKey string set *jwt.Token at *gin.Context
func Auth(secret, signMethod, contextKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// header we want
		//header.Set("Authorization", fmt.Sprintf("Bearer %v", token.Value()))
		var auth, authhead string
		authhead = c.Request.Header.Get("Authorization")

		if strings.HasPrefix(authhead, bearer_) {
			auth = strings.Split(authhead, bearer_)[1]
		} else {
			c.AbortWithError(401, errors.New("invalid jwt header"))
		}

		token, err := jwt.Parse(auth, func(t *jwt.Token) (interface{}, error) {
			// Check the signing method
			if t.Method.Alg() != signMethod {
				return nil, fmt.Errorf("unexpected jwt signing method %s", t.Header["alg"])
			}
			return []byte(secret), nil
		})

		if err != nil {
			c.AbortWithError(401, err)
		}

		if contextKey != "" {
			// do have context Key name
			c.Set(contextKey, token)
		}
	}
}
