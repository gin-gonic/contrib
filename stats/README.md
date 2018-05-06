# Gin's middleware for request stats

[![Build Status](https://travis-ci.org/semihalev/gin-stats.svg)](https://travis-ci.org/semihalev/gin-stats)
[![codecov](https://codecov.io/gh/semihalev/gin-stats/branch/master/graph/badge.svg)](https://codecov.io/gh/semihalev/gin-stats)
[![Go Report Card](https://goreportcard.com/badge/github.com/semihalev/gin-stats)](https://goreportcard.com/report/github.com/semihalev/gin-stats)
[![GoDoc](https://godoc.org/github.com/semihalev/gin-stats?status.svg)](https://godoc.org/github.com/semihalev/gin-stats)

Lightweight Gin's middleware for request metrics

## Usage

### Start using it

Download and install it:

```sh
$ go get github.com/semihalev/gin-stats
```

```go
import "github.com/semihalev/gin-stats"
```

### Example usage:

```go
package main

import (
	"fmt"
	"time"

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

```