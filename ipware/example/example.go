package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rounds/go-gin-contrib/ipware"
)

func main() {
	r := gin.New()
	r.Use(ipware.UpdateAddr())
	r.GET("/", func(c *gin.Context) {
		c.String(200, "real ip: "+fmt.Sprint(c.Request.RemoteAddr))
	})
	r.Run(":8080")
}
