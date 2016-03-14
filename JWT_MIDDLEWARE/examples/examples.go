package main

import (
	"fmt"
	"github.com/akshaykumar12527/contrib/JWT_MIDDLEWARE"
	"github.com/gin-gonic/gin"
	"time"
)

func main() {
	router := gin.Default()
	jwt_middleware := JWT_MIDDLEWARE.JWTMiddleware{
		Key:     []byte("thisisreallybigandsecurekey"),
		Realm:   "jwt auth",
		Timeout: time.Hour * 24 * 10,
		Authenticator: func(userId string, password string) bool {
			fmt.Println("into jwt_middleware")
			return true
		},
	}

	// This handler will match /user/john but will not match neither /user/ or /user
	router.GET("/user", jwt_middleware.MiddlewareFunc(), func(c *gin.Context) {
		fmt.Println("/user")
		c.JSON(200, gin.H{"message": "success"})
		return
	})

	// However, this one will match /user/john/ and also /user/john/send
	// If no other routers match /user/john, it will redirect to /user/john/
	router.PUT("/user", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})

	})

	router.Run(":8080")
}
