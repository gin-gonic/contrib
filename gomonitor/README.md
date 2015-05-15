# gomonitor
Gin middleware for exposing metrics with
[go-monitor](https://github.com/mcuadros/go-monitor). It supports
custom monitors

## Examples

### Custom
```go
package main

import (
	"github.com/gin-gonic/contrib/gomonitor"
	"github.com/gin-gonic/gin"
	"gopkg.in/mcuadros/go-monitor.v1"
	"gopkg.in/mcuadros/go-monitor.v1/aspects"
)

func setupMetrics(aspect *CustomAspect) {
	m := monitor.NewMonitor(":9000")
	m.AddAspect(aspect)
	m.Start()
}

type CustomAspect struct {
	Count int
}

func (a *CustomAspect) GetStats() interface{} {
	a.Count++
	return a.Count
}

func (a *CustomAspect) Name() string {
	return "Custom"
}

func (a *CustomAspect) InRoot() bool {
	return false
}
func main() {
	anAspect := &CustomAspect{3}
	asps := []aspects.Aspect{anAspect}
	router := gin.New()
    // curl http://localhost:9000/
    // curl http://localhost:9000/Custom
	router.Use(gomonitor.Metrics(9000, asps))
    //..
    router.Run(":8080")
}
```
