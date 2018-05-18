// Package ginrus provides log handling using logrus package.
//
// Based on github.com/stephenmuss/ginerus but adds more options.
package ginrus

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type loggerEntryWithFields interface {
	WithFields(fields logrus.Fields) *logrus.Entry
}

// Ginrus returns a gin.HandlerFunc (middleware) that logs requests using logrus.
//
// Requests with errors are logged using logrus.Error().
// Requests without errors are logged using logrus.Info().
//
// It receives:
//   1. A time package format string (e.g. time.RFC3339).
//   2. A boolean stating whether to use UTC time zone or local.
//   3. Optionally, paths to skip logging for.
func Ginrus(logger loggerEntryWithFields, timeFormat string, utc bool, notlogged ...string) gin.HandlerFunc {
	var skip map[string]struct{}
	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)
		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}
	return func(c *gin.Context) {
		start := time.Now()
		// some evil middlewares modify this values
		path := c.Request.URL.Path
		c.Next()

		// log only when path is not being skippd
		if _, ok := skip[path]; !ok {
			end := time.Now()
			latency := end.Sub(start)
			if utc {
				end = end.UTC()
			}

			entry := logger.WithFields(logrus.Fields{
				"status":     c.Writer.Status(),
				"method":     c.Request.Method,
				"path":       path,
				"ip":         c.ClientIP(),
				"latency":    latency,
				"user-agent": c.Request.UserAgent(),
				"time":       end.Format(timeFormat),
			})

			if len(c.Errors) > 0 {
				// Append error field if this is an erroneous request.
				entry.Error(c.Errors.String())
			} else {
				entry.Info()
			}
		}
	}
}
