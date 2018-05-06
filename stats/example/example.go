package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/semihalev/gin-stats"
)

func main() {
	r := gin.Default()
	r.Use(stats.RequestStats())

	r.GET("/stats", func(c *gin.Context) {
		c.JSON(http.StatusOK, stats.Report())
	})

	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
