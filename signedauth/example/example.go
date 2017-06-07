package main

import (
	"github.com/gin-gonic/contrib/signedauth"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	mgr := StrictSHA1Manager{Prefix: "SAUTH", Key: "contextKey", Secret: "super-secret-password", Value: nil}
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
