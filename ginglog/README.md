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

func setupYourAppLogging() {
	go func() {
		for _ = range time.Tick(1 * time.Second) {
			glog.Flush()
		}
	}()
}

func main() {
	parseFlags() // send flags to glog
	setupYourAppLogging()
	router := gin.New()
    // curl http://localhost:8080/
    // curl http://localhost:8080/ping
    router.Use(ginglog.Logger(2 * time.Second))
    //..
    router.Run(":8080")
}
```
