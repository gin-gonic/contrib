package main

import (
	"github.com/ChristopherRabotin/gin-contrib/signedauth"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	mgr := StrictSHA1Manager{prefix: "SAUTH", key: "contextKey", secret: "super-secret-password", value: nil}
	router := gin.Default()
	router.Use(signedauth.SignatureAuth(mgr))
	router.POST("/test/", func(c *gin.Context) {
			c.String(http.StatusOK, "Success.")
		})
	router.PUT("/test/", func(c *gin.Context) {
			c.String(http.StatusOK, "Success.")
		})
	router.Run("localhost:31337")
}
