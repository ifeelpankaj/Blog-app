package app

import (
	"blog_app/config"
	"blog_app/middleware"
	"blog_app/utils/logger"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"

	"github.com/danielkov/gin-helmet/ginhelmet"
)

func limitRequestBody(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		c.Next()
	}
}
func StartApp() {
	config.LoadConfig()
	logger.InitLogger(config.AppConfig.Env)

	// Middleware
	r := gin.Default()
	r.Use(ginhelmet.Default())
	r.Use(middleware.CORS())

	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.Use(static.Serve("/static", static.LocalFile("./static", true)))

	r.Use(limitRequestBody(50 << 20))

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Server is running on "+config.AppConfig.Env+" mode")
	})
	r.POST("/upload", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Upload received"})
	})

	if config.AppConfig.Env == "prod" {
		r.SetTrustedProxies([]string{"your-prod-proxy-ip"})
	}

	// Start server
	port := config.AppConfig.Port
	r.Run(":" + port)

}
