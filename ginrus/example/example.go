package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/contrib/ginrus"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()

	// Add a ginrus middleware, which logs all requests to stdout in RFC3339 format, with UTC support.
	r.Use(ginrus.Ginrus(logrus.StandardLogger(), false, time.RFC3339, true))

	// Add another ginrus middleware which logs errors to stderr, in local time.
	logger := logrus.New()
	logger.SetOutput(os.Stderr)
	r.Use(ginrus.Ginrus(logger, true, time.RFC3339, false))

	// Ping
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong "+fmt.Sprint(time.Now().Unix()))
	})

	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
