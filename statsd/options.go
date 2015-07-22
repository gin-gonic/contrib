// Here be option setting functions for the statsd middleware.

package statsd

const (
	DefaultResponseTimeEnabled = true
	DefaultThroughputEnabled   = true
	DefaultStatusCodeEnabled   = true
	DefaultSuccessEnabled      = true
	DefaultErrorEnabled        = true

	DefaultResponseTimeBucket = "request.response_time"
	DefaultThroughputBucket   = "request.throughput"
	DefaultStatusCodeBucket   = "request.status_code"
	DefaultSuccessBucket      = "request.success"
	DefaultErrorBucket        = "request.error.default"
)

// OptionFunc is a configuration function, used in Statsd().
type OptionFunc func(*configuredClient)

func SetResponseTime(enabled bool) OptionFunc {
	return func(s *configuredClient) {
		s.ResponseTimeEnabled = enabled
	}
}
func SetThroughput(enabled bool) OptionFunc {
	return func(s *configuredClient) {
		s.ThroughputEnabled = enabled
	}
}
func SetStatusCode(enabled bool) OptionFunc {
	return func(s *configuredClient) {
		s.StatusCodeEnabled = enabled
	}
}
func SetSuccess(enabled bool) OptionFunc {
	return func(s *configuredClient) {
		s.SuccessEnabled = enabled
	}
}
func SetError(enabled bool) OptionFunc {
	return func(s *configuredClient) {
		s.ErrorEnabled = enabled
	}
}

func SetResponseTimeBucket(name string) OptionFunc {
	return func(s *configuredClient) {
		s.ResponseTimeBucket = name
	}
}
func SetThroughputBucket(name string) OptionFunc {
	return func(s *configuredClient) {
		s.ThroughputBucket = name
	}
}
func SetStatusCodeBucket(name string) OptionFunc {
	return func(s *configuredClient) {
		s.StatusCodeBucket = name
	}
}
func SetSuccessBucket(name string) OptionFunc {
	return func(s *configuredClient) {
		s.SuccessBucket = name
	}
}
func SetErrorBucket(name string) OptionFunc {
	return func(s *configuredClient) {
		s.ErrorBucket = name
	}
}
