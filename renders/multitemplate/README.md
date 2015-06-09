This is a custom HTML render to support multi templates, ie. more than one `*template.Template`.



```go
package main

import (
    "html/template"

    "github.com/gin-gonic/gin"
    "github.com/gin-gonic/contrib/renders/multitemplate"
)

func main() {
    router := gin.Default()
    router.HTMLRender = createMyRender()
    router.GET("/", func(c *gin.Context) {
        c.HTML(200, "index", data)
    })
    router.Run(":8080")
}

func createMyRender() multitemplate.Render {
    r := multitemplate.New()
    r.AddFromFiles("index", "base.html", "base.html")
    r.AddFromFiles("article", "base.html", "article.html")
    r.AddFromFiles("login", "base.html", "login.html")
    r.AddFromFiles("dashboard", "base.html", "dashboard.html")

    return r
}
```