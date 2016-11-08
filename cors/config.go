package cors

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

const toLower = 'a' - 'A'

type converter func(string) string

type settings struct {
	allowAllOrigins   bool
	allowedOriginFunc func(string) bool
	allowedOrigins    []string
	allowedMethods    []string
	allowedHeaders    []string
	exposedHeaders    []string
	normalHeaders     http.Header
	preflightHeaders  http.Header
}

func newSettings(c Config) *settings {
	if err := c.Validate(); err != nil {
		panic(err.Error())
	}
	return &settings{
		allowedOriginFunc: c.AllowOriginFunc,
		allowAllOrigins:   c.AllowAllOrigins,
		allowedOrigins:    c.AllowedOrigins,
		allowedMethods:    distinct(convert(c.AllowedMethods, strings.ToUpper)),
		allowedHeaders:    distinct(convert(c.AllowedHeaders, http.CanonicalHeaderKey)),
		normalHeaders:     generateNormalHeaders(c),
		preflightHeaders:  generatePreflightHeaders(c),
	}
}

func (c *settings) validateOrigin(origin string) (string, bool) {
	if c.allowAllOrigins {
		return "*", true
	}
	if c.allowedOriginFunc != nil {
		return origin, c.allowedOriginFunc(origin)
	}
	for _, value := range c.allowedOrigins {
		if value == origin {
			return origin, true
		}
	}
	return "", false
}

func (c *settings) validateMethod(method string) bool {
	if len(c.allowedMethods) == 0 {
		return false
	}
	method = strings.ToUpper(method)
	if method == "OPTIONS" {
		// Always allow preflight requests
		return true
	}
	for _, m := range c.allowedMethods {
		if m == method {
			return true
		}
	}
	return false
}

func (c *settings) validateHeaders(requestHeaders []string) bool {
	if len(requestHeaders) == 0 {
		return true
	}
	for _, header := range requestHeaders {
		header = http.CanonicalHeaderKey(header)
		found := false
		for _, h := range c.allowedHeaders {
			if h == header {
				found = true
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func generateNormalHeaders(c Config) http.Header {
	headers := make(http.Header)
	if c.AllowCredentials {
		headers.Set("Access-Control-Allow-Credentials", "true")
	}
	if len(c.ExposedHeaders) > 0 {
		headers.Set("Access-Control-Expose-Headers", strings.Join(convert(c.ExposedHeaders, http.CanonicalHeaderKey), ", "))
	}
	return headers
}

func generatePreflightHeaders(c Config) http.Header {
	headers := make(http.Header)
	if c.AllowCredentials {
		headers.Set("Access-Control-Allow-Credentials", "true")
	}
	if len(c.AllowedMethods) > 0 {
		headers.Set("Access-Control-Allow-Methods", strings.Join(convert(c.AllowedMethods, strings.ToUpper), ", "))
	}
	if len(c.AllowedHeaders) > 0 {
		headers.Set("Access-Control-Allow-Headers", strings.Join(convert(c.AllowedHeaders, http.CanonicalHeaderKey), ", "))
	}
	if c.MaxAge > time.Duration(0) {
		headers.Set("Access-Control-Max-Age", strconv.FormatInt(int64(c.MaxAge/time.Second), 10))
	}
	return headers
}

func distinct(s []string) []string {
	m := map[string]bool{}
	for _, v := range s {
		if _, seen := m[v]; !seen {
			s[len(m)] = v
			m[v] = true
		}
	}
	return s[:len(m)]
}

func parse(content string) []string {
	if len(content) == 0 {
		return nil
	}
	parts := strings.Split(content, ",")
	for i := 0; i < len(parts); i++ {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func convert(s []string, c converter) []string {
	var out []string
	for _, i := range s {
		out = append(out, c(i))
	}
	return out
}

func parseHeaders(headers string) []string {
	l := len(headers)
	h := make([]byte, 0, l)
	upper := true
	// Estimate the number headers in order to allocate the right splice size
	t := 0
	for i := 0; i < l; i++ {
		if headers[i] == ',' {
			t++
		}
	}
	headerList := make([]string, 0, t)
	for i := 0; i < l; i++ {
		b := headers[i]
		if b >= 'a' && b <= 'z' {
			if upper {
				h = append(h, b-toLower)
			} else {
				h = append(h, b)
			}
		} else if b >= 'A' && b <= 'Z' {
			if !upper {
				h = append(h, b+toLower)
			} else {
				h = append(h, b)
			}
		} else if b == '-' || b == '_' || (b >= '0' && b <= '9') {
			h = append(h, b)
		}

		if b == ' ' || b == ',' || i == l-1 {
			if len(h) > 0 {
				// Flush the found header
				headerList = append(headerList, string(h))
				h = h[:0]
				upper = true
			}
		} else {
			upper = b == '-' || b == '_'
		}
	}
	return headerList
}
