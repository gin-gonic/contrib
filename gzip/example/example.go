package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/rounds/go-gin-contrib/gzip"
)

func main() {
	r := gin.Default()
	r.Use(gzip.OutputFilter(nil))
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong "+fmt.Sprint(time.Now().Unix()))
	})

	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
