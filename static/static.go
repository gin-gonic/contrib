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
func Serve(directories ...interface{}) gin.HandlerFunc {
	fileservers := []http.Handler{}

	for i := 0; i < len(directories); i++ {
		directory, err := filepath.Abs(directories[i].(string))
		if err != nil {
			panic(err)
		}
		fileservers = append(fileservers, http.FileServer(http.Dir(directory)))
	}

	return func(c *gin.Context) {
		for i := 0; i < len(directories); i++ {
			directory := directories[i].(string)
			p := path.Join(directory, c.Request.URL.Path)
			if exists(p) {
				fileservers[i].ServeHTTP(c.Writer, c.Request)
				c.Abort()
				break
			}
		}
	}
}
