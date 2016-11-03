package main

import (
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-gonic/contrib/static"
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"strings"
)

type binaryFileSystem struct {
	fs http.FileSystem
}

func (b *binaryFileSystem) Open(name string) (http.File, error) {
	return b.fs.Open(name)
}

func (b *binaryFileSystem) Exists(prefix string, filepath string) bool {

	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		if _, err := b.fs.Open(p); err != nil {
			return false
		}
		return true
	}
	return false
}

func BinaryFileSystem(root string) *binaryFileSystem {
	fs := &assetfs.AssetFS{Asset, AssetDir, root}
	return &binaryFileSystem{
		fs,
	}
}

// Usage
// $ go-bindata data/
// $ go build && ./bindata
//
func main() {
	r := gin.Default()

	r.Use(static.Serve("/static", BinaryFileSystem("data")))
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "test")
	})
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
