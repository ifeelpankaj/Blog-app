package middleware

import (
	"blog_app/config"

	"blog_app/utils/logger"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CORS() gin.HandlerFunc {
	config.LoadConfig()
	logger.InitLogger(config.AppConfig.Env)
	allowedOrigin := config.AppConfig.AllowedOrigin

	return cors.New(cors.Config{
		AllowOrigins:     []string{allowedOrigin},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowCredentials: true,
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		MaxAge:           12 * time.Hour,
		AllowOriginFunc: func(origin string) bool {
			if origin == "" {
				return true // Allow non-origin requests (e.g., mobile apps)
			}
			if origin != allowedOrigin {
				logger.Warn("CORS policy violation attempt from origin: %s", zap.String("origin", origin))
				return false
			}
			return true
		},
	})
}
