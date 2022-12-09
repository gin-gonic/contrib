package statsd

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// configuredClient wraps a StatsD Client with configuration settings.
type configuredClient struct {
	Client Client

	// On/off switch for each metric.
	ResponseTimeEnabled,
	ThroughputEnabled,
	StatusCodeEnabled,
	SuccessEnabled,
	ErrorEnabled bool

	// Bucket names for each metric.
	ResponseTimeBucket,
	ThroughputBucket,
	StatusCodeBucket,
	SuccessBucket,
	ErrorBucket string
}

func newConfiguredClient(client Client) *configuredClient {
	return &configuredClient{
		Client: client,

		ResponseTimeEnabled: DefaultResponseTimeEnabled,
		ThroughputEnabled:   DefaultThroughputEnabled,
		StatusCodeEnabled:   DefaultStatusCodeEnabled,
		SuccessEnabled:      DefaultSuccessEnabled,
		ErrorEnabled:        DefaultErrorEnabled,

		ResponseTimeBucket: DefaultResponseTimeBucket,
		ThroughputBucket:   DefaultThroughputBucket,
		StatusCodeBucket:   DefaultStatusCodeBucket,
		SuccessBucket:      DefaultSuccessBucket,
		ErrorBucket:        DefaultErrorBucket,
	}
}

// IncrThroughput increments the received requests bucket.
func (c *configuredClient) IncrThroughput(handler string) {
	if c.ThroughputEnabled {
		c.Client.Incr(join(handler, c.ThroughputBucket), 1)
	}
}

// IncrStatusCode increments the context's response status code bucket.
func (c *configuredClient) IncrStatusCode(status int, handler string) {
	if c.StatusCodeEnabled {
		c.Client.Incr(join(handler, c.StatusCodeBucket, strconv.Itoa(status)), 1)
	}
}

// IncrSuccess increments the success bucket
// if no errors were attached to the context.
func (c *configuredClient) IncrSuccess(errors []*gin.Error, handler string) {
	if c.SuccessEnabled && len(errors) == 0 {
		c.Client.Incr(join(handler, c.SuccessBucket), 1)
	}
}

// IncrError increments the error bucket for each attached error.
// If a gin.Error.Meta implements the Bucket interface,
// its associated bucket will be incremented instead of the default error bucket.
// See the Bucket interface for more info.
//
// NOTE that if at least one error was attached to the context,
// calling IncrSucess does nothing.
func (c *configuredClient) IncrError(errors []*gin.Error, handler string) {
	if c.ErrorEnabled {
		for _, err := range errors {
			// If the gin.Error.Meta implements the Bucket interface,
			// increment its specific bucket by its given increment amount.
			// Otherwise, increment the default error bucket with the default amount.
			b, ok := err.Meta.(Bucket)
			if !ok {
				c.Client.Incr(join(handler, c.ErrorBucket), 1)
				continue
			}
			c.Client.Incr(join(handler, b.BucketName()), 1)
		}
	}
}

// Timing calculates time taken to process the request.
func (c *configuredClient) Timing(start time.Time, handler string) {
	if c.ResponseTimeEnabled {
		c.Client.Timing(join(handler, c.ResponseTimeBucket),
			// Convert to milliseconds.
			time.Now().Sub(start).Nanoseconds()/time.Millisecond.Nanoseconds())
	}
}

// join is a helper function that joins strings to a valid StatsD bucket name,
// separated by a dot '.'.
//
// It's main use is to have cleaner code in this library.
func join(strs ...string) string { return strings.Join(strs, ".") }
