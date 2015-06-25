package resourcer

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

type User struct {
	Id       int
	Name     string
	Password string
}

type GormResourcer struct {
	ResourceGorm
}

type GormCollectioner struct {
	CollectionGorm
}

func (_ *GormCollectioner) ResourceClass() (reflect.Type, string) {
	return reflect.TypeOf(User{}), "user"
}

func (_ *GormResourcer) ResourceClass() (reflect.Type, string) {
	return reflect.TypeOf(User{}), "user"
}

func TestResourcer(t *testing.T) {
	router := gin.New()

	called := false
	mockdb, err := sqlmock.New()
	assert.Nil(t, err)
	db, err := gorm.Open("mysql", mockdb)
	assert.Nil(t, err)
	group := router.Group("/users")

	group.Use(ResourcerGorm(&GormResourcer{}, &db))
	group.GET("/:user_id", func(c *gin.Context) {
		called = true
		_, ok := c.Get("gin.contrib.resource")
		assert.True(t, ok)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/1", nil)
	sqlmock.ExpectQuery("\\*").WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "password"}).
		AddRow(1, "username", "password"))
	router.ServeHTTP(w, req)
	assert.True(t, called)
	assert.Equal(t, w.Code, 200)
}

func TestCollectioner(t *testing.T) {
	router := gin.New()
	called := false

	mockdb, err := sqlmock.New()
	assert.Nil(t, err)
	db, err := gorm.Open("mysql", mockdb)
	assert.Nil(t, err)
	group := router.Group("/users")

	group.Use(ResourcerGorm(&GormCollectioner{}, &db))
	group.GET("/photos", func(c *gin.Context) {
		called = true
		_, ok := c.Get("gin.contrib.resource")
		assert.True(t, ok)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/photos", nil)

	sqlmock.ExpectQuery("\\*").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "password"}).
		AddRow(1, "username", "password"))
	router.ServeHTTP(w, req)
	assert.True(t, called)
	assert.Equal(t, w.Code, 200)
}

func TestResourcer404(t *testing.T) {
	router := gin.New()

	mockdb, err := sqlmock.New()
	assert.Nil(t, err)
	db, err := gorm.Open("mysql", mockdb)
	assert.Nil(t, err)
	group := router.Group("/users")

	group.Use(ResourcerGorm(&GormResourcer{}, &db))
	group.GET("/:user_id", func(c *gin.Context) {})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/notaninteger", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, w.Code, 404)
	gin.SetMode(gin.DebugMode)
	router.ServeHTTP(w, req)
	assert.Equal(t, w.Code, 404)
}
