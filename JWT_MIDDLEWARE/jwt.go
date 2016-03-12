// Package jwt provides Json-Web-Token authentication for the go-json-rest framework
package JWT_MIDDLEWARE

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	// "net/http"
	"strings"
	"time"
)

// JWTMiddleware provides a Json-Web-Token authentication implementation. On failure, a 401 HTTP response
// is returned. On success, the wrapped middleware is called, and the userId is made available as
// request.Env["userID"].(string).
// Users can get a token by posting a json request to LoginHandler. The token then needs to be passed in
// the Authentication header. Example: Authorization:Bearer XXX_TOKEN_XXX#!/usr/bin/env
type JWTMiddleware struct {
	// Realm name to display to the user. Required.
	Realm string

	// signing algorithm - possible values are HS256, HS384, HS512
	// Optional, default is HS256.
	SigningAlgorithm string

	// Secret key used for signing. Required.
	Key []byte

	// Duration that a jwt token is valid. Optional, defaults to one hour.
	Timeout time.Duration

	// This field allows clients to refresh their token until MaxRefresh has passed.
	// Note that clients can refresh their token in the last moment of MaxRefresh.
	// This means that the maximum validity timespan for a token is MaxRefresh + Timeout.
	// Optional, defaults to 0 meaning not refreshable.
	MaxRefresh time.Duration

	// Callback function that should perform the authentication of the user based on userId and
	// password. Must return true on success, false on failure. Required.
	Authenticator func(userId string, password string) bool

	// Callback function that should perform the authorization of the authenticated user. Called
	// only after an authentication success. Must return true on success, false on failure.
	// Optional, default to success.
	Authorizator func(userId string, c *gin.Context) bool

	// Callback function that will be called during login.
	// Using this function it is possible to add additional payload data to the webtoken.
	// The data is then made available during requests via request.Env["JWT_PAYLOAD"].
	// Note that the payload is not encrypted.
	// The attributes mentioned on jwt.io can't be used as keys for the map.
	// Optional, by default no additional data will be set.
	PayloadFunc func(userId string) map[string]interface{}
}

// MiddlewareFunc makes JWTMiddleware implement the Middleware interface.
func (mw *JWTMiddleware) MiddlewareFunc() gin.HandlerFunc {
	// fmt.Println("MiddlewareFunc")
	if mw.Realm == "" {
		log.Fatal("Realm is required")
	}
	if mw.SigningAlgorithm == "" {
		mw.SigningAlgorithm = "HS256"
	}
	if mw.Key == nil {
		log.Fatal("Key required")
	}
	if mw.Timeout == 0 {
		mw.Timeout = time.Hour
	}
	if mw.Authenticator == nil {
		log.Fatal("Authenticator is required")
	}
	if mw.Authorizator == nil {
		mw.Authorizator = func(userId string, c *gin.Context) bool {
			return true
		}
	}

	return func(c *gin.Context) {

		mw.middlewareImpl(c)
	}
}

func (mw *JWTMiddleware) middlewareImpl(c *gin.Context) {

	fmt.Println("middlewareImpl")

	token, err := mw.parseToken(c)

	if err != nil {
		mw.unauthorized(c)
		return
	}

	uid := token.Claims["id"].(string)

	// fmt.Println(loudshoutUsers)
	c.Set("userID", uid)
	c.Next()
}

// Helper function to extract the JWT claims
func ExtractClaims(c *gin.Context) map[string]interface{} {
	fmt.Println("ExtractClaims")
	if val, _ := c.Get("JWT_PAYLOAD"); val == nil {
		empty_claims := make(map[string]interface{})
		return empty_claims
	}
	jwt_claims, _ := c.Get("JWT_PAYLOAD")
	return jwt_claims.(map[string]interface{})
}

// Handler that clients can use to get a jwt token.
// Payload needs to be json in the form of {"username": "USERNAME", "password": "PASSWORD"}.
// Reply will be of the form {"token": "TOKEN"}.
func (mw *JWTMiddleware) TokenGenerator(c *gin.Context, userID string) string {
	fmt.Println("LoginHandler")
	// mw.SigningAlgorithm = "HS256"
	fmt.Println(mw.SigningAlgorithm, string(mw.Key))

	token := jwt.New(jwt.GetSigningMethod(mw.SigningAlgorithm))
	// fmt.Println(token)
	if mw.PayloadFunc != nil {
		for key, value := range mw.PayloadFunc(userID) {
			token.Claims[key] = value
		}
	}
	// uid := userID.String()
	// fmt.Println("userID=",uid)
	token.Claims["id"] = userID
	token.Claims["exp"] = time.Now().Add(mw.Timeout).Unix()

	tokenString, err := token.SignedString(mw.Key)

	if err != nil {
		mw.unauthorized(c)
		return "null"
	}
	return tokenString
}

func (mw *JWTMiddleware) parseToken(c *gin.Context) (*jwt.Token, error) {
	authHeader := c.Request.Header.Get("Authorization")
	// fmt.Println(authHeader)
	if authHeader == "" {
		return nil, errors.New("Auth header empty")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return nil, errors.New("Invalid auth header")
	}

	return jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod(mw.SigningAlgorithm) != token.Method {
			return nil, errors.New("Invalid signing algorithm")
		}
		return mw.Key, nil
	})
}

type token struct {
	Token string `json:"token"`
}

func (mw *JWTMiddleware) unauthorized(c *gin.Context) {
	fmt.Println("unauthorized")
	c.Request.Header.Set("WWW-Authenticate", "JWT realm="+mw.Realm)

	c.JSON(401, gin.H{"userMessege": "Not Authorized"})
	// rest.Error(writer, "Not Authorized", http.StatusUnauthorized)
}
