package gzip

import "compress/gzip"

// Compression levels
// These constants reference gzip package, so that code that imports
// "contrib/compress" does not also have to import "compress/gzip".
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

// Defaults
var defaults = struct {
	CompressionLevel     int
	ForceCompression     bool
	CompressContentTypes []string
}{
	CompressionLevel: DefaultCompression,
	ForceCompression: false,
	CompressContentTypes: []string{
		"text/plain",
		"text/html",
		"text/javascript",
		"application/javascript",
		"application/x-javascript",
		"application/json",
		"application/octet-stream",
	},
}
