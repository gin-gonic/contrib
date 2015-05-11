package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rounds/gin-gonic/contrib/gzip"
)

func main() {
	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	// Ping
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong "+fmt.Sprint(time.Now().Unix()))
	})

	// Echo
	r.POST("/echo", func(c *gin.Context) {
		req, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.Fail(500, err)
		}
		// echo the request text
		requestText := string(req)
		c.String(200, requestText)
	})

	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
