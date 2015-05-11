package gzip

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const (
	testRequest  = "Gzip Test Request"
	testResponse = "Gzip Test Response "
)

func newServer(useGzip bool) *gin.Engine {
	r := gin.Default()
	if useGzip {
		r.Use(Gzip(DefaultCompression))
	}
	r.GET("/", func(c *gin.Context) {
		c.Writer.Header().Set(headerContentLength, strconv.Itoa(len(testResponse)))
		c.String(200, testResponse)
	})
	r.POST("/", func(c *gin.Context) {
		req, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.Fail(500, err)
		}
		// echo the request text
		requestText := string(req)
		c.String(200, requestText)
	})
	return r
}

func TestGetGzip(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add(headerAcceptEncoding, encodingGzip)

	w := httptest.NewRecorder()

	r := newServer(true)
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusOK)

	enc := w.Header().Get(headerContentEncoding)
	assert.Equal(t, encodingGzip, enc)

	enc = w.Header().Get(headerVary)
	assert.Equal(t, headerAcceptEncoding, enc)

	length := w.Header().Get(headerContentLength)
	assert.Equal(t, "", length)
	assert.NotEqual(t, len(testResponse), w.Body.Len())

	gr, err := gzip.NewReader(w.Body)
	assert.NoError(t, err)
	defer gr.Close()

	body, _ := ioutil.ReadAll(gr)
	assert.Equal(t, testResponse, string(body))
}

func TestGetNoGzip(t *testing.T) {
	// check that all works without the middleware
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add(headerAcceptEncoding, encodingGzip)
	w := httptest.NewRecorder()
	r := newServer(false)
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusOK)
	enc := w.Header().Get(headerContentEncoding)
	assert.Equal(t, "", enc)

	length := w.Header().Get(headerContentLength)
	assert.Equal(t, strconv.Itoa(len(testResponse)), length)
	assert.Equal(t, testResponse, w.Body.String())

	// add the middleware and check that all is well
	req, _ = http.NewRequest("GET", "/", nil)
	w = httptest.NewRecorder()
	r = newServer(true)
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusOK)
	enc = w.Header().Get(headerContentEncoding)
	assert.Equal(t, "", enc)

	length = w.Header().Get(headerContentLength)
	assert.Equal(t, strconv.Itoa(len(testResponse)), length)
	assert.Equal(t, testResponse, w.Body.String())
}

func TestPostGzip(t *testing.T) {
	var b bytes.Buffer
	g := gzip.NewWriter(&b)
	g.Write([]byte(testRequest))
	g.Close()

	req, _ := http.NewRequest("POST", "/", &b)
	req.Header.Add(headerContentEncoding, encodingGzip)

	w := httptest.NewRecorder()

	r := newServer(true)
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusOK)

	length := w.Header().Get(headerContentLength)
	assert.Equal(t, "", length)

	body, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, testRequest, string(body))
}

func TestPostNoGzip(t *testing.T) {
	// check not compressed post with the middeware enabled
	req, _ := http.NewRequest("POST", "/", strings.NewReader(testRequest))
	w := httptest.NewRecorder()
	r := newServer(true)
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusOK)
	assert.Equal(t, testRequest, w.Body.String())

	// check not compressed post without the middeware
	req, _ = http.NewRequest("POST", "/", strings.NewReader(testRequest))
	w = httptest.NewRecorder()
	r = newServer(false)
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusOK)
	//length := w.Header().Get(headerContentLength)
	//assert.Equal(t, strconv.Itoa(len(testRequest)), length)
	assert.Equal(t, testRequest, w.Body.String())
}
