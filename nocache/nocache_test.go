package nocache

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNoCache(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/test", nil)
	g := gin.New()
	g.Use(NoCache())
	g.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"test": "test",
		})
	})
	g.ServeHTTP(w, r)

	for k, v := range noCacheHeaders {
		t.Run(k, func(t *testing.T) {
			require.Equal(t, w.Header().Get(k), v)
		})
	}
}
