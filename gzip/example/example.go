package main

import (
	"fmt"
	"github.com/gin-gonic/contrib/gzip"
	"gopkg.in/gin-gonic/gin.v1"
	"time"
)

func main() {
	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong "+fmt.Sprint(time.Now().Unix()))
	})

	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
