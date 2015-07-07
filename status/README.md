# Status

Status middleware for go gonic. 

Feel free to add more data

# Example

```
package main

import (
    "github.com/gin-gonic/contrib/status"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    statusMw := &status.StatusMiddleware{}
    r.Use(statusMw.Status())

    r.GET("/.status", func(c *gin.Context) {
        c.JSON(200, statusMw.GetStatus())
    })

    r.Run("localhost:8080")
}
```

when sending GET request to /.status path the output will look like

```
{
  "Pid": 24269,
  "UpTime": "42.475291007s",
  "Time": "2015-07-07 03:40:51.923920701 +0300 EEST",
  "TimeUnix": 1436229651,
  "PerStatus": {
    "200": {
      "RequestCount": 3,
      "TotalResponseTimesPerRequest": "1µs",
      "AverageResponseTimesPerRequest": "0"
    },
    "401": {
      "RequestCount": 3,
      "TotalResponseTimesPerRequest": "17µs",
      "AverageResponseTimesPerRequest": "6µs"
    },
    "404": {
      "RequestCount": 2,
      "TotalResponseTimesPerRequest": "0",
      "AverageResponseTimesPerRequest": "0"
    },
    "501": {
      "RequestCount": 3,
      "TotalResponseTimesPerRequest": "146µs",
      "AverageResponseTimesPerRequest": "49µs"
    }
  },
  "TotalCount": 11,
  "TotalResponseTime": "163.242µs",
  "AverageResponseTime": "14.84µs"
}
```