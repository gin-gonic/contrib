package gzip

import (
	"compress/gzip"
	"strings"

	"github.com/gin-gonic/gin"
)

// Compression levels
const (
	BestCompression    = gzip.BestCompression
	BestSpeed          = gzip.BestSpeed
	DefaultCompression = gzip.DefaultCompression
	NoCompression      = gzip.NoCompression
)

// Internal constants
const (
	encodingGzip = "gzip"

	headerAcceptEncoding  = "Accept-Encoding"
	headerContentEncoding = "Content-Encoding"
	headerContentLength   = "Content-Length"
	headerContentType     = "Content-Type"
	headerVary            = "Vary"
)

// Our custom writer
type gzipWriter struct {
	gin.ResponseWriter
	gzwriter *gzip.Writer
}

func (g *gzipWriter) Write(data []byte) (int, error) {
	return g.gzwriter.Write(data)
}

// Gzip is a Gin middleware to handle HTTP request/response compression with gzip.
// Requests having compressed body must set the Content-Encoding header.
// To receive a compressed response, the request must have the Accept-Encoding header.
func Gzip(level int) gin.HandlerFunc {

	return func(c *gin.Context) {
		req := c.Request

		// Decode the request payload if applicable.
		if strings.Contains(req.Header.Get(headerContentEncoding), encodingGzip) {
			gzreader, err := gzip.NewReader(req.Body)
			if err != nil {
				c.Next()
				return
			}
			req.Body = gzreader
		}

		// See if the response can be encoded.
		if !strings.Contains(req.Header.Get(headerAcceptEncoding), encodingGzip) {
			c.Next()
			return
		}

		// Init gzip writer engine.
		gz, err := gzip.NewWriterLevel(c.Writer, level)
		if err != nil {
			c.Next()
			return
		}

		// Set compression-related headers.
		headers := c.Writer.Header()
		headers.Set(headerContentEncoding, encodingGzip)
		headers.Set(headerVary, headerAcceptEncoding)

		// Save the original writer and replace it with our own.
		writer := c.Writer
		c.Writer = &gzipWriter{c.Writer, gz}

		// Render response.
		c.Next()

		// Remove content-length header because it's now incorrect.
		writer.Header().Del(headerContentLength)

		// Clean up and return to the original writer.
		gz.Close()
		c.Writer = writer
	}
}
