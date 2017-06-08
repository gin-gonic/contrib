package cors

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Config struct {
	AbortOnError    bool
	AllowAllOrigins bool

	// AllowedOrigins is a list of origins a cross-domain request can be executed from.
	// If the special "*" value is present in the list, all origins will be allowed.
	// Default value is ["*"]
	AllowedOrigins []string

	// AllowOriginFunc is a custom function to validate the origin. It take the origin
	// as argument and returns true if allowed or false otherwise. If this option is
	// set, the content of AllowedOrigins is ignored.
	AllowOriginFunc func(origin string) bool

	// AllowedMethods is a list of methods the client is allowed to use with
	// cross-domain requests. Default value is simple methods (GET and POST)
	AllowedMethods []string

	// AllowedHeaders is list of non simple headers the client is allowed to use with
	// cross-domain requests.
	// If the special "*" value is present in the list, all headers will be allowed.
	// Default value is [] but "Origin" is always appended to the list.
	AllowedHeaders []string

	// ExposedHeaders indicates which headers are safe to expose to the API of a CORS
	// API specification
	ExposedHeaders []string

	// AllowCredentials indicates whether the request can include user credentials like
	// cookies, HTTP authentication or client side SSL certificates.
	AllowCredentials bool

	// MaxAge indicates how long (in seconds) the results of a preflight request
	// can be cached
	MaxAge time.Duration
}

func (c *Config) AddAllowedMethods(methods ...string) {
	c.AllowedMethods = append(c.AllowedMethods, methods...)
}

func (c *Config) AddAllowedHeaders(headers ...string) {
	c.AllowedHeaders = append(c.AllowedHeaders, headers...)
}

func (c *Config) AddExposedHeaders(headers ...string) {
	c.ExposedHeaders = append(c.ExposedHeaders, headers...)
}

func (c Config) Validate() error {
	if c.AllowAllOrigins && (c.AllowOriginFunc != nil || len(c.AllowedOrigins) > 0) {
		return errors.New("conflict settings: all origins are allowed. AllowOriginFunc or AllowedOrigins is not needed")
	}
	if !c.AllowAllOrigins && c.AllowOriginFunc == nil && len(c.AllowedOrigins) == 0 {
		return errors.New("conflict settings: all origins disabled")
	}
	if c.AllowOriginFunc != nil && len(c.AllowedOrigins) > 0 {
		return errors.New("conflict settings: if a allow origin func is provided, AllowedOrigins is not needed")
	}
	for _, origin := range c.AllowedOrigins {
		if !strings.HasPrefix(origin, "http://") && !strings.HasPrefix(origin, "https://") {
			return errors.New("bad origin: origins must include http:// or https://")
		}
	}
	return nil
}

var defaultConfig = Config{
	AbortOnError:    false,
	AllowAllOrigins: true,
	AllowedMethods:  []string{"GET", "POST", "PUT", "PATCH", "HEAD"},
	AllowedHeaders:  []string{"Content-Type"},
	//ExposedHeaders:   "",
	AllowCredentials: false,
	MaxAge:           12 * time.Hour,
}

func DefaultConfig() Config {
	cp := defaultConfig
	return cp
}

func Default() gin.HandlerFunc {
	return New(defaultConfig)
}

func New(config Config) gin.HandlerFunc {
	s := newSettings(config)

	// Algorithm based in http://www.html5rocks.com/static/images/cors_server_flowchart.png
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if len(origin) == 0 {
			return
		}
		origin, valid := s.validateOrigin(origin)
		if valid {
			c.Header("Access-Control-Allow-Origin", origin)
			if c.Request.Method == "OPTIONS" {
				valid = handlePreflight(c, s)
				c.Status(200)
				return
			} else {
				valid = handleNormal(c, s)
			}
		}

		if !valid {
			if config.AbortOnError {
				c.AbortWithStatus(http.StatusForbidden)
			}
			return
		}
		c.Next()
	}
}

func handlePreflight(c *gin.Context, s *settings) bool {
	// Always set Vary headers
	// see https://github.com/rs/cors/issues/10,
	// https://github.com/rs/cors/commit/dbdca4d95feaa7511a46e6f1efb3b3aa505bc43f#commitcomment-12352001
	c.Writer.Header().Add("Vary", "Origin")
	c.Writer.Header().Add("Vary", "Access-Control-Request-Method")
	c.Writer.Header().Add("Vary", "Access-Control-Request-Headers")

	if !s.validateMethod(c.Request.Header.Get("Access-Control-Request-Method")) {
		return false
	}
	if !s.validateHeaders(parseHeaders(c.Request.Header.Get("Access-Control-Request-Headers"))) {
		return false
	}
	for key, value := range s.preflightHeaders {
		c.Writer.Header()[key] = value
	}

	return true
}

func handleNormal(c *gin.Context, s *settings) bool {
	c.Writer.Header().Add("Vary", "Origin")

	for key, value := range s.normalHeaders {
		c.Writer.Header()[key] = value
	}
	return true
}
