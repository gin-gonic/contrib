package stats

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestStats(t *testing.T) {
	router := gin.New()
	router.Use(RequestStats())

	router.GET("/stats", func(c *gin.Context) {
		c.JSON(http.StatusOK, Report())
	})

	w := performRequest(router, "GET", "/stats")
	assert.Equal(t, w.Code, 200)
	assert.Equal(t, w.Body.String(), "{}")
}
