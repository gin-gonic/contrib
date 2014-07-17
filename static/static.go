package static

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

func exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

// Static returns a middleware handler that serves static files in the given directory.
func Serve(directory string) gin.HandlerFunc {
	directory, err := filepath.Abs(directory)
	if err != nil {
		panic(err)
	}
	fileserver := http.FileServer(http.Dir(directory))

	return func(c *gin.Context) {

		p := path.Join(directory, c.Request.URL.Path)
		if exists(p) {
			fileserver.ServeHTTP(c.Writer, c.Request)
			c.Abort(-1)
		}
	}
}
