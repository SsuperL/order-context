package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck ...
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, "OK")
}

// StartHealthCheckHTTPServer ...
func StartHealthCheckHTTPServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}
