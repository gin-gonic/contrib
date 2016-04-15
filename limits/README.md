# Various Limits Middleware for gin

## RateLimiter

### Example:

    package main

    import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/contrib/ratelimit"
    )


    func handler(ctx *gin.Context) {
            val := ctx.PostForm("b")
            if len(ctx.Errors) > 0 {
                    return
            }
            ctx.String(http.StatusOK, "got %s\n", val)
    }

    func main() {
            rtr := gin.Default()
            rtr.Use(ratelimit.RateLimiter(10))
            rtr.POST("/", handler)
            rtr.Run(":8080")
    }
