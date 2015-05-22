package main

import (
  "github.com/gin-gonic/gin"
  "github.com/ChristopherRabotin/gin-contrib/signedauth"
  "net/http"
)


func main(){
  mgr := StrictSHA1Manager{prefix: "SAUTH", key: "contextKey", secret: "super-secret-password", value: nil}
		router := gin.Default()
		router.Use(signedauth.SignatureAuth(mgr))
		methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
		for _, meth := range methods {
			router.Handle(meth, "/test/", []gin.HandlerFunc{func(c *gin.Context) {
				c.String(http.StatusOK, "Success.")
			}})
		}
		router.Run("localhost:31337")
}
