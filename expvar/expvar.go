package expvar

import (
	"expvar"
	"fmt"

	"gopkg.in/gin-gonic/gin.v1"
)

func Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		w := c.Writer
		c.Header("Content-Type", "application/json; charset=utf-8")
		w.Write([]byte("{\n"))
		first := true
		expvar.Do(func(kv expvar.KeyValue) {
			if !first {
				w.Write([]byte(",\n"))
			}
			first = false
			fmt.Fprintf(w, "%q: %s", kv.Key, kv.Value)
		})
		w.Write([]byte("\n}\n"))
		c.AbortWithStatus(200)
	}
}
