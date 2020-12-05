// Package statsd provides reporting HTTP performance to a StatsD agent.
package statsd

import (
	"time"

	"github.com/gin-gonic/gin"
)

// The Client interface is implemented by a StatsD client.
// The middleware receives a Client object and uses it to report metrics.
//
// You can use any StatsD client you wish,
// as long as it implements this interface.
//
// I wrote this middleware to be used with the most popular package:
// github.com/quipo/statsd
type Client interface {
	Incr(stat string, count int64) error
	Timing(stat string, delta int64) error
}

// The Bucket interface is implemented by gin.Error.Meta types
// who wish their associated error to increment a different bucket
// INSTEAD of the default error bucket.
//
// E.g. Consider two error types E and F.
// E implements the Bucket interface while F does not.
// Upon an E error, the E.BucketName() bucket will be incremented.
// Upon an F error, the generic error bucket will be incremented.
type Bucket interface {
	BucketName() string // Bucket name to increment.
}

// Statsd returns a gin.HandleFunc (middleware) that sends metrics to the given StatsD Client.
//
// The following metrics are collected and reported to the given StatsD agent:
//   1. Request throughput count.
//   2. Response status code count.
//   3. Response time.
//   4. Successful response count.
//      A response is considered "successful" if no errors were attached to its context.
//   5. Context error count.
//      If three errors were attached to the context, three increments will occur.
//
// Bucket names have a default value, which can be set by the option functions.
//
// Furthermore, if a gin.Error.Meta type implements the Bucket interface,
// its associated bucket will be incremented INSTEAD of the default error bucket.
// See the Bucket interface for more info.
func Statsd(client Client, options ...OptionFunc) gin.HandlerFunc {
	// Initialize and configure the client and set options if given.
	cc := newConfiguredClient(client)
	for _, option := range options {
		option(cc)
	}

	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		handler := c.HandlerName()
		cc.IncrThroughput(handler)
		cc.IncrStatusCode(c.Writer.Status(), handler)
		cc.IncrSuccess(c.Errors, handler)
		cc.IncrError(c.Errors, handler)
		cc.Timing(start, handler)
	}
}
