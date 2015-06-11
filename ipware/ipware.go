// Package ipware updates gin.Conext.Request.RemoteAddr with the client's real
// IP address on a best attempt basis.
//
// Based on https://github.com/sebest/xff
package ipware

import (
	"github.com/gin-gonic/gin"
	"github.com/sebest/xff"
)

// UpdateAddr returns a gin.HandlerFunc (middleware) that
// updates gin.Context.Request.RemoteAddr with the client's real IP address,
// on a best attempt basis.
//
// It parses Forwarded HTTP extension headers (RFC 7239).
//
// See the following article for more information: http://tools.ietf.org/html/rfc7239
func UpdateAddr() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.RemoteAddr = xff.GetRemoteAddr(c.Request)
		c.Next()
	}
}
