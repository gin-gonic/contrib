package gomonitor

import (
	"fmt"

	"github.com/gin-gonic/gin"
	mon "gopkg.in/mcuadros/go-monitor.v1"
	"gopkg.in/mcuadros/go-monitor.v1/aspects"
)

var monitor *mon.Monitor

func Metrics(port int, asp []aspects.Aspect) gin.HandlerFunc {
	monitor := mon.NewMonitor(fmt.Sprintf(":%d", port))
	for _, aspect := range asp {
		monitor.AddAspect(aspect)
	}

	go monitor.Start()
	return func(c *gin.Context) {
		c.Next()
	}
}
