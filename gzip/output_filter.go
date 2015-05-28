package gzip

import (
	"compress/gzip"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// OutputOptions for the middleware.
type OutputOptions struct {
	CompressionLevel     int
	ForceCompression     bool
	CompressContentTypes []string
}

var (
	outputOptions = &OutputOptions{
		defaults.CompressionLevel,
		defaults.ForceCompression,
		defaults.CompressContentTypes,
	}
)

// Our custom writer
type gzipWriter struct {
	gin.ResponseWriter
	gzwriter *gzip.Writer
}

// Write implements the Writer interface.
func (g *gzipWriter) Write(data []byte) (int, error) {
	// TODO:
	// If we still don't know the content type, we can check
	// g.ResponseWriter.Header() or try to detect it using http.DetectContentType.
	return g.gzwriter.Write(data)
}

func mimeTypeSupported(mtype string) bool {
	for _, ctype := range outputOptions.CompressContentTypes {
		if strings.HasPrefix(mtype, ctype) {
			return true
		}
	}
	return false
}

// OutputFilter is a Gin middleware to handle HTTP response compression.
// To receive a compressed response, the request must have the Accept-Encoding header.
func OutputFilter(opts *OutputOptions) gin.HandlerFunc {

	if opts != nil {
		outputOptions = opts
	}

	return func(c *gin.Context) {
		req := c.Request

		// See if there's an Accept-Encoding header.
		if !strings.Contains(req.Header.Get(headerAcceptEncoding), encodingGzip) {
			// Abort if compression is mantatory.
			if outputOptions.ForceCompression {
				c.AbortWithStatus(http.StatusNotAcceptable)
				return
			}
			// Continue normal flow, without compression.
			c.Next()
			return
		}

		// See if extension is given and we support this mime type.
		extension := filepath.Ext(req.URL.Path)
		if extension != "" {
			mtype := mime.TypeByExtension(extension)
			if !mimeTypeSupported(mtype) {
				// Continue normal flow, without compression.
				c.Next()
				return
			}
		}

		// Init gzip writer engine.
		gz, err := gzip.NewWriterLevel(c.Writer, outputOptions.CompressionLevel)
		if err != nil {
			c.Error(err)
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
		defer func() {
			// Remove content-length header because it's now incorrect.
			writer.Header().Del(headerContentLength)
			// Clean up and return to the original writer.
			gz.Close()
			c.Writer = writer
		}()

		// Render response.
		c.Next()
	}
}
