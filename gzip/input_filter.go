package gzip

import (
	"compress/gzip"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// InputOptions for the middleware.
type InputOptions struct {
	ForceCompression bool
}

var (
	inputOptions = &InputOptions{defaults.ForceCompression}
)

// InputFilter is a Gin middleware to handle HTTP request compression.
// Requests having compressed body must set the Content-Encoding header.
func InputFilter(opts *InputOptions) gin.HandlerFunc {

	if opts != nil {
		inputOptions = opts
	}

	return func(c *gin.Context) {
		req := c.Request

		// See if there's a Content-Encoding header.
		if !strings.Contains(req.Header.Get(headerContentEncoding), encodingGzip) {
			// Abort if compression is mantatory.
			if inputOptions.ForceCompression {
				c.AbortWithStatus(http.StatusUnsupportedMediaType)
				return
			}
			// Continue normal flow, without compression.
			c.Next()
			return
		}

		// Init Gzip reader.
		gzreader, err := gzip.NewReader(req.Body)
		if err != nil {
			c.Error(err)
			c.Next()
			return
		}

		// Alright, replace the reader and continue.
		req.Body = gzreader
		c.Next()
	}
}
