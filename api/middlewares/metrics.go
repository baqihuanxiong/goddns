package middlewares

import (
	"github.com/gin-gonic/gin"
	"goddns/metrics"
	"strconv"
	"time"
)

func Metrics(service metrics.Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		// start time
		start := time.Now()

		// process request
		c.Next()

		// save metrics
		statusCode := strconv.Itoa(c.Writer.Status())
		duration := float64(time.Since(start).Milliseconds())
		service.CollectHttp(c.HandlerName(), c.Request.Method, statusCode, duration)
	}
}
