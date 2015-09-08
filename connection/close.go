package connection

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

// CloseAfterPrefix will ensure that the connection is explicitly closed by the
// server if the route matches urlPrefix.
func CloseAfterPrefix(urlPrefix string) gin.HandlerFunc {
	if urlPrefix == "" {
		log.Println("setting closeConnectionOn prefix to \"/\"")
		urlPrefix = "/"
	}

	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, urlPrefix) {
			c.Request.Header.Set("Connection", "close")
		}
		c.Next()
	}
}

// CloseAfterAll will close all http requests after the route has finished
// serving
func CloseAfterAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Header.Set("Connection", "close")
		c.Next()
	}
}
