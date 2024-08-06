package main

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	gin.DisableConsoleColor()

	router := gin.Default()

	// Trust only addresses configured with the PROXIES envvar.
	if proxies := strings.Split(os.Getenv("PROXIES"), ","); proxies[0] != "" {
		if err := router.SetTrustedProxies(proxies); err != nil {
			panic(err)
		}
	}

	// Install middleware.
	router.Use(rateLimitMiddleware()) // Called 1st.
	router.Use(corsMiddleware())      // Called 2nd.

	registerEndpoints(router)

	return router
}

func main() {
	const address = "127.0.0.1:9172"

	router := setupRouter()
	err := router.Run(address)
	if err != nil {
		panic(err)
	}
}
