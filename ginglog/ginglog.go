package ginglog

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

func setupLogging(duration time.Duration) {
	go func() {
		for _ = range time.Tick(duration) {
			glog.Flush()
		}
	}()
}

var (
	green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow  = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset   = string([]byte{27, 91, 48, 109})
)

func ErrorLogger() gin.HandlerFunc {
	return ErrorLoggerT(gin.ErrorTypeAll)
}

func ErrorLoggerT(typ uint32) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		errs := c.Errors.ByType(typ)
		if len(errs) > 0 {
			// -1 status code = do not change current one
			c.JSON(-1, c.Errors)
		}
	}
}

func Logger(duration time.Duration) gin.HandlerFunc {
	setupLogging(duration)
	return func(c *gin.Context) {
		t := time.Now()

		// process request
		c.Next()

		latency := time.Since(t)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		statusColor := colorForStatus(statusCode)
		methodColor := colorForMethod(method)

		glog.Infof("[GIN] |%s %3d %s| %12v | %s |%s  %s %-7s %s\n%s",
			statusColor, statusCode, reset,
			latency,
			clientIP,
			methodColor, reset, method,
			c.Request.URL.Path,
			c.Errors.String(),
		)

	}
}

func colorForStatus(code int) string {
	switch {
	case code >= 200 && code <= 299:
		return green
	case code >= 300 && code <= 399:
		return white
	case code >= 400 && code <= 499:
		return yellow
	default:
		return red
	}
}

func colorForMethod(method string) string {
	switch {
	case method == "GET":
		return blue
	case method == "POST":
		return cyan
	case method == "PUT":
		return yellow
	case method == "DELETE":
		return red
	case method == "PATCH":
		return green
	case method == "HEAD":
		return magenta
	case method == "OPTIONS":
		return white
	default:
		return reset
	}
}
