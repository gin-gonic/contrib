package status

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type StatusMiddleware struct {
	lock                    sync.RWMutex
	start                   time.Time
	pid                     int
	responseCounts          map[string]int
	responseTimesPerRequest map[string]float64
	totalResponseTime       time.Time
}

func (mw *StatusMiddleware) Status() gin.HandlerFunc {
	mw.start = time.Now()
	mw.pid = os.Getpid()
	mw.responseCounts = map[string]int{}
	mw.responseTimesPerRequest = map[string]float64{}
	mw.totalResponseTime = time.Time{}

	return func(c *gin.Context) {

		start := time.Now()

		// Process Request
		c.Next()

		end := time.Now()
		responseTime := end.Sub(start)
		statusCode := fmt.Sprintf("%d", c.Writer.Status())

		mw.lock.Lock()
		mw.responseCounts[statusCode]++
		mw.responseTimesPerRequest[statusCode] += responseTime.Seconds()
		mw.totalResponseTime = mw.totalResponseTime.Add(responseTime)
		mw.lock.Unlock()
	}
}

type dataPerRequest struct {
	RequestCount                   int
	TotalResponseTimesPerRequest   string
	AverageResponseTimesPerRequest string
}

type Status struct {
	Pid                 int
	UpTime              string
	Time                string
	TimeUnix            int64
	PerStatus           map[string]dataPerRequest
	TotalCount          int
	TotalResponseTime   string
	AverageResponseTime string
}

func (mw *StatusMiddleware) GetStatus() *Status {

	mw.lock.RLock()

	now := time.Now()

	uptime := now.Sub(mw.start)

	totalCount := 0
	for _, count := range mw.responseCounts {
		totalCount += count
	}

	perStatus := map[string]dataPerRequest{}

	for status, _ := range mw.responseCounts {
		var data dataPerRequest

		averageAsString := fmt.Sprintf("%.6fs", mw.responseTimesPerRequest[status]/float64(mw.responseCounts[status]))
		totalAsString := fmt.Sprintf("%.6fs", mw.responseTimesPerRequest[status])

		averageResponseTimesPerRequest, _ := time.ParseDuration(averageAsString)
		totalResponseTimesPerRequest, _ := time.ParseDuration(totalAsString)

		data.AverageResponseTimesPerRequest = averageResponseTimesPerRequest.String()
		data.TotalResponseTimesPerRequest = totalResponseTimesPerRequest.String()
		data.RequestCount = mw.responseCounts[status]

		perStatus[status] = data
	}

	totalResponseTime := mw.totalResponseTime.Sub(time.Time{})

	averageResponseTime := time.Duration(0)
	if totalCount > 0 {
		avgNs := int64(totalResponseTime) / int64(totalCount)
		averageResponseTime = time.Duration(avgNs)
	}

	status := &Status{
		Pid:                 mw.pid,
		UpTime:              uptime.String(),
		Time:                now.String(),
		TimeUnix:            now.Unix(),
		PerStatus:           perStatus,
		TotalCount:          totalCount,
		TotalResponseTime:   totalResponseTime.String(),
		AverageResponseTime: averageResponseTime.String(),
	}

	mw.lock.RUnlock()

	return status
}
