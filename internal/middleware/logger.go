package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
)

func LoggerMiddleware(logger *logrus.Entry) gin.HandlerFunc {

	return func(c *gin.Context) {

		start := time.Now()
		c.Next()
		duration := time.Since(start).Milliseconds()

		path := c.Request.URL.Path
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()
		referer := c.Request.Referer()

		entry := logger.WithFields(logrus.Fields{
			"statusCode": statusCode,
			"duration":   duration,
			"clientIP":   clientIP,
			"method":     c.Request.Method,
			"path":       path,
			"referer":    referer,
			"userAgent":  userAgent,
		})

		entry.Infof(
			"[GIN] %s - %s \"%s\" %d (%dms) ",
			clientIP,
			c.Request.Method,
			path,
			statusCode,
			duration,
		)
	}
}
