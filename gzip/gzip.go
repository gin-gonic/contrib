package gzip

import (
	"compress/gzip"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	BestCompression    = gzip.BestCompression
	BestSpeed          = gzip.BestSpeed
	DefaultCompression = gzip.DefaultCompression
	NoCompression      = gzip.NoCompression
)

func Gzip(level int) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !clientAcceptGzip(c.Request) {
			return
		}

		gz, err := gzip.NewWriterLevel(c.Writer, level)
		if err != nil {
			return
		}
		defer gz.Close()
		c.Header("Content-Encoding", "gzip")
		c.Header("Vary", "Accept-Encoding")

		c.Writer = &gzipWriter{c.Writer, gz}
		c.Next()
		c.Header("Content-Length", "")
	}
}

func clientAcceptGzip(req *http.Request) bool {
	return strings.Contains(req.Header.Get("Accept-Encoding"), "gzip")
}

type gzipWriter struct {
	gin.ResponseWriter
	gzwriter *gzip.Writer
}

func (g *gzipWriter) Write(data []byte) (int, error) {
	return g.gzwriter.Write(data)
}
