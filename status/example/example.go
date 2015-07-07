package main

import (
	"net/http"

	"github.com/gin-gonic/contrib/status"
	"github.com/gin-gonic/gin"
)

var statusMw *status.StatusMiddleware

func main() {
	r := gin.Default()

	statusMw = &status.StatusMiddleware{}
	r.Use(statusMw.Status())

	r.GET("/.status", Status)
	r.GET("/fail", func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"message": "Hello World"})
	})

	r.GET("/noauth", func(c *gin.Context) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Hello World"})
	})

	r.GET("/abort", func(c *gin.Context) {
		c.Abort()
	})

	r.GET("/abortwithstat", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusUnauthorized)
	})

	r.Run("localhost:8080")
}

func Status(c *gin.Context) {
	c.JSON(200, statusMw.GetStatus())
}
