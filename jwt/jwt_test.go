package jwt

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const secretKey = "MyTestSigningKey"

func TestAuth(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.GET("/private", Auth(secretKey, jwt.SigningMethodHS256), privateHandler)

	response := makeRequest(router, "GET", "/private", "")

	// No token
	if response.Code != 401 {
		t.Errorf("No token.  Expected 401, got %d", response.Code)
	}

	// Empty token
	response = makeRequest(router, "GET", "/private", "Authorization: Bearer ")
	if response.Code != 401 {
		t.Errorf("Empty token: Expected 401, got %d", response.Code)
	}

	// Token signed with wrong key
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["auth"] = true
	tokenString, _ := token.SignedString([]byte("WrongKey"))
	response = makeRequest(router, "GET", "/private", fmt.Sprintf("Bearer %s", tokenString))
	if response.Code != 401 {
		t.Errorf("Empty token: Expected 401, got %d", response.Code)
	}

	// Token signed with right key
	token = jwt.New(jwt.SigningMethodHS256)
	token.Claims["auth"] = true
	tokenString, _ = token.SignedString([]byte("MyTestSigningKey"))
	response = makeRequest(router, "GET", "/private", fmt.Sprintf("Bearer %s", tokenString))
	if response.Code != 200 {
		t.Errorf("Correct token: Expected 200, got %d", response.Code)
	}
	response.Flush()

	var body map[string]interface{}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Errorf("Failed to unmarshal response: %s", err)
		return
	}

	claims, ok := body["claims"].(map[string]interface{})
	if !ok {
		t.Errorf("Claims missing from body")
		return
	}

	auth, ok := claims["auth"].(bool)
	if !ok {
		t.Errorf("Missing auth claim")
		return
	}

	if !auth {
		t.Errorf("Auth was somehow valid but false")
	}
}

func privateHandler(c *gin.Context) {
	var claims map[string]interface{}
	if cl, exists := c.Get("claims"); exists {
		var ok bool
		claims, ok = cl.(map[string]interface{})
		if !ok {
			c.AbortWithError(401, fmt.Errorf("missing claims"))
		}
	}

	c.JSON(200, gin.H{
		"claims": claims,
	})
}

func makeRequest(r http.Handler, method, path, authHeader string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	if authHeader != "" {
		req.Header.Add("Authorization", authHeader)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
