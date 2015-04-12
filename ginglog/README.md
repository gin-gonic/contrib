# ginglog
Gin middleware for Logging with
[glog](https://github.com/golang/glog). It is meant as drop in
replacement for the default logging used in gin.

## Examples

### Custom

Start your webapp to log to STDERR and /tmp:

    % ./webapp -log_dir="/tmp" -logtostderr

```go
package main

import (
    "time"

    "github.com/golang/glog" // needed if you use glog for your app logging
    "github.com/gin-gonic/contrib/ginglog"
	"github.com/gin-gonic/gin"
)

func main() {
	parseFlags() // send flags to glog
	router := gin.New()
    router.Use(ginglog.Logger(3 * time.Second))
    //..
    glog.Info("bootstrapped application")
    router.Run(":8080")
}
```
