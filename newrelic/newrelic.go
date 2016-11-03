package newrelic

import (
	"time"

	"gopkg.in/gin-gonic/gin.v1"
	metrics "github.com/yvasiyarov/go-metrics"
	"github.com/yvasiyarov/gorelic"
)

var agent *gorelic.Agent

func NewRelic(license string, appname string, verbose bool) gin.HandlerFunc {
	agent = gorelic.NewAgent()
	agent.NewrelicLicense = license

	agent.HTTPTimer = metrics.NewTimer()
	agent.CollectHTTPStat = true
	agent.Verbose = verbose

	agent.NewrelicName = appname
	agent.Run()

	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		agent.HTTPTimer.UpdateSince(startTime)
	}
}
