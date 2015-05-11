// Package ginrus provides log handling using logrus package.
//
// Based on github.com/stephenmuss/ginerus but adds more options.
package ginrus

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

// Ginrus returns a gin.HandlerFunc (middleware) that logs requests using logrus.
//
// Errors are logged using logrus.Error() function. The rest are using logrus.Info().
//
// It receives:
//   1. A boolean stating whether to log only errors.
//   2. A time package format string (e.g. time.RFC3339).
//   3. A boolean stating whether to use UTC time zone or local.
func Ginrus(logger *logrus.Logger, errorsOnly bool, timeFormat string, utc bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now().UTC()
		path := c.Request.URL.Path

		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		errors := c.Errors.String()
		userAgent := c.Request.UserAgent()

		if !errorsOnly || (errorsOnly && len(c.Errors) > 0) {
			entry := logger.WithFields(logrus.Fields{
				"status":     statusCode,
				"method":     method,
				"error":      errors,
				"ip":         clientIP,
				"latency":    latency,
				"path":       path,
				"time":       end.Format(timeFormat),
				"user-agent": userAgent,
			})

			if len(c.Errors) > 0 {
				entry.Error()
			} else {
				entry.Info()
			}
		}
	}
}
