package main

import (
	"fmt"
	"time"

	"github.com/uber-go/zap"
	"github.com/yezooz/contrib/ginzap"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()

	log := zap.New(
		zap.NewTextEncoder(),
		zap.DebugLevel,
	)

	// Add a ginzap middleware, which:
	//   - Logs all requests, like a combined access and error log.
	//   - Logs to stdout.
	//   - RFC3339 with UTC time format.
	r.Use(ginzap.Ginzap(log, time.RFC3339, true))

	// Example ping request.
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong "+fmt.Sprint(time.Now().Unix()))
	})

	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
