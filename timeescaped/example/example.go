package main

import (
	"fmt"
	"github.com/gin-gonic/contrib/timeescaped"
	"github.com/gin-gonic/gin"
	"time"
)

func main() {
	r := gin.Default()

	r.Use(timeescaped.TimeEscaped())

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong "+fmt.Sprint(time.Now().Unix()))
	})
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
