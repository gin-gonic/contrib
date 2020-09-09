# NoCache

NoCache is a simple piece of middleware that sets a number of HTTP headers to prevent a router (or subrouter) from being cached by an upstream proxy and/or client.

## Usage
```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/nocache"
)

func main() {
    g := gin.New()
    g.Use(nocache.NoCache())
    g.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "result": "It will not be cached",
        })
    })
}
```