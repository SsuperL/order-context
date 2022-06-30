package common

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// LogMiddleWare ...
func LogMiddleWare() gin.HandlerFunc {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.DebugLevel)
	// logger.SetFormatter(&logrus.JSONFormatter{})

	return func(c *gin.Context) {
		//Start time
		startTime := time.Now()

		//Process request
		c.Next()

		//End time
		endTime := time.Now()

		//Execution time
		latencyTime := endTime.Sub(startTime)

		//Request method
		reqMethod := c.Request.Method

		//Request routing
		reqUri := c.Request.RequestURI

		// status code
		statusCode := c.Writer.Status()

		// request IP
		clientIP := c.ClientIP()
		logger.WithFields(logrus.Fields{
			"status_code":  statusCode,
			"latency_time": latencyTime,
			"req_method":   reqMethod,
			"req_uri":      reqUri,
			"client_ip":    clientIP,
		}).Info()

	}

}

// SideCodeHanler ...
func SideCodeHanler() gin.HandlerFunc {

	return func(c *gin.Context) {
		siteCode := c.Request.Header.Get("site-code")
		if siteCode == "" {
			c.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{"error": "site-code is required"})
			return
		}

		// c.Request.WithContext()

	}

}
