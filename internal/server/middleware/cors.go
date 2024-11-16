package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
	corsConfig.AllowHeaders = []string{
		"Authorization", "Content-Type", "Origin", "session-key", "Api-Token",
	}
	return cors.New(corsConfig)
}
