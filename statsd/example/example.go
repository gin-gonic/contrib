package main

import (
	"fmt"
	"time"

	middleware "github.com/gin-gonic/contrib/statsd"
	"github.com/gin-gonic/gin"
	"github.com/quipo/statsd"
)

func main() {
	// Initialize a StatsD agent of your choice,
	// which implements the Client interface.
	client := statsd.NewStatsdClient("127.0.0.1", "my.app.")
	bufferedClient := statsd.NewStatsdBuffer(5*time.Second, client)

	r := gin.New()

	// Add a StatsD middleware with a custom throughput bucket name.
	// See the other option functions for additional settings.
	r.Use(middleware.Statsd(bufferedClient,
		middleware.SetThroughputBucket("request.received")))

	// Example ping request.
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong "+fmt.Sprint(time.Now().Unix()))
	})

	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
