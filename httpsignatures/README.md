#  httpsignatures

Gin middleware base on [HTTP Signatures](https://tools.ietf.org/html/draft-cavage-http-signatures).

## Example
``` go
package main

import (
	"github.com/gin-gonic/contrib/httpsignatures"
	"github.com/gin-gonic/gin"
)

func main() {
	// Define algorithm
	hmacsha256 := &httpsignatures.HmacSha256{}
	hmacsha512 := &httpsignatures.HmacSha512{}
	// Init define secret params
	readKeyID := httpsignatures.KeyID("read")
	writeKeyID := httpsignatures.KeyID("write")
	secrets := httpsignatures.Secrects{
		readKeyID: &httpsignatures.Secret{
			Key:       "1234",
			Algorithm: hmacsha256, // You could using other algo with interface Crypto
		},
		writeKeyID: &httpsignatures.Secret{
			Key:       "5678",
			Algorithm: hmacsha512,
		},
	}
	// Define permission list
	writePermissions := httpsignatures.Permission{writeKeyID}
	readPermissions := httpsignatures.Permission{readKeyID, writeKeyID}

	// Init server
	r := gin.Default()

	// Init require params of middleware
	requiredHeaders := []string{"(request-target)", "date", "digest"}
	dateValidator := httpsignatures.NewDateValidator()
	auth := httpsignatures.NewAuthenticator(secrets, requiredHeaders, dateValidator)

	// Group required read permission
	readPermissionHandler := auth.Authenticated(readPermissions)
	read := r.Group("/read")
	read.Use(readPermissionHandler)
	read.GET("/a", a)
	read.POST("/b", b)
	read.POST("/c", c)

	// Group required write permission
	writePermissionHandler := auth.Authenticated(writePermissions)
	write := r.Group("/write")
	write.Use(writePermissionHandler)
	write.GET("/x", x)
	write.POST("/y", y)
	write.POST("/z", z)

	r.Run(":8080")
}

```