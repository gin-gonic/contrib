package static

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type ServeFileSystem interface {
	http.FileSystem
	Exists(prefix string, path string) bool
}

func existsFile(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

type localFileSystem struct {
	fs      http.FileSystem
	root    string
	indexes bool
}

func LocalFile(root string, indexes bool) *localFileSystem {
	root, err := filepath.Abs(root)
	if err != nil {
		panic(err)
	}

	fs := http.Dir(root)
	return &localFileSystem{
		fs,
		root,
		indexes,
	}
}

func (l *localFileSystem) Open(name string) (http.File, error) {
	f, err := l.fs.Open(name)
	if err != nil {
		return nil, err
	}

	if l.indexes {
		return f, err
	} else {
		return neuteredReaddirFile{f}, nil
	}
}

func (l *localFileSystem) Exists(prefix string, filepath string) bool {

	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		p = path.Join(l.root, p)
		return existsFile(p)
	}
	return false
}

type neuteredReaddirFile struct {
	http.File
}

func (f neuteredReaddirFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

// Static returns a middleware handler that serves static files in the given directory.
func Serve(prefix string, fs ServeFileSystem) gin.HandlerFunc {
	var fileserver http.Handler

	if prefix != "" {
		fileserver = http.StripPrefix(prefix, http.FileServer(fs))
	} else {
		fileserver = http.FileServer(fs)
	}

	return func(c *gin.Context) {

		if fs.Exists(prefix, c.Request.URL.Path) {
			fileserver.ServeHTTP(c.Writer, c.Request)
		} else {
			c.Next()
		}
	}
}
