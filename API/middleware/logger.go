package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		requestID := uuid.New().String()
		c.Writer.Header().Set("X-Request-ID", requestID)

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		gin.DefaultWriter.Write([]byte(
			requestID + " | " +
				c.Request.Method + " " +
				path + " | " +
				statusString(status) + " | " +
				latency.String() + "\n",
		))
	}
}

func statusString(code int) string {
	return fmt.Sprintf("%d", code)
}
