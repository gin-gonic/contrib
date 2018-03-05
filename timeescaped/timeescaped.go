package timeescaped

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// TimeEscaped calc the time cost in per request
func TimeEscaped() gin.HandlerFunc {

	return func(c *gin.Context) {
		start := time.Now()
		c.Next() // next handler func
		fmt.Printf("path:%s meth:%s cost time:%s\n", c.Request.URL.Path, c.Request.Method, time.Now().Sub(start).String())
	}
}
