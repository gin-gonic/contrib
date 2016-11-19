package pprof

import (
	"net/http"
	"net/http/pprof"

	"github.com/gin-gonic/gin"
)

// Options provides potential route registration configuration options
type Options struct {
	// RoutePrefix is an optional path prefix
	RoutePrefix string
}

// Register the standard HandlerFuncs from the net/http/pprof package with
// the provided gin.Engine
func Register(r *gin.Engine, opts *Options) {
	prefix := routePrefix(opts)
	r.GET(prefix+"/debug/pprof/block", pprofHandler(pprof.Index))
	r.GET(prefix+"/debug/pprof/heap", pprofHandler(pprof.Index))
	r.GET(prefix+"/debug/pprof/profile", pprofHandler(pprof.Profile))
	r.POST(prefix+"/debug/pprof/symbol", pprofHandler(pprof.Symbol))
	r.GET(prefix+"/debug/pprof/symbol", pprofHandler(pprof.Symbol))
	r.GET(prefix+"/debug/pprof/trace", pprofHandler(pprof.Trace))
}

func pprofHandler(h http.HandlerFunc) gin.HandlerFunc {
	handler := http.HandlerFunc(h)
	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}

func routePrefix(opts *Options) string {
	if opts == nil {
		return ""
	}
	return opts.RoutePrefix
}
