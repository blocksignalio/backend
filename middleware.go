package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func corsMiddleware() gin.HandlerFunc {
	origins := strings.Split(os.Getenv("ORIGINS"), ",")

	if origins[0] == "" {
		return cors.Default()
	}

	config := cors.DefaultConfig()
	config.AllowOrigins = origins
	return cors.New(config)
}

func rateLimitMiddleware() gin.HandlerFunc {
	const message = "You can only make one request per second."

	limiter := rate.NewLimiter(10, 30)
	return func(c *gin.Context) {
		if limiter.Allow() {
			c.Next()
		} else {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    codeClientError,
				"message": message,
				"data":    nil,
			})
			c.Abort()
		}
		c.Next()
	}
}
