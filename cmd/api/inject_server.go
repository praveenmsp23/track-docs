package main

import (
	"github.com/google/wire"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/praveenmsp23/trackdocs/handler/api"
	"github.com/praveenmsp23/trackdocs/handler/health"
	"github.com/praveenmsp23/trackdocs/pkg/config"
	"github.com/praveenmsp23/trackdocs/pkg/token"
)

// wire set for loading the server.
var serverSet = wire.NewSet(
	health.NewHealth,
	api.NewApi,
	provideRouter,
)

func provideRouter(cfg *config.Config, manager *token.Manager, health *health.Health, api *api.Api) *gin.Engine {
	if cfg.Env == config.ApplicationEnvLocal {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()
	engine.RedirectTrailingSlash = false
	engine.MaxMultipartMemory = 2 << 20 // 2 MiB
	engine.Use(gin.Recovery())
	engine.Use(CORSMiddleware(cfg))
	health.Routes(engine.Group("/health"))
	api.Routes(engine.Group("/api"))
	engine.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error_code": http.StatusNotFound, "error_message": "page not found"})
	})
	engine.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"success": false, "error_code": http.StatusMethodNotAllowed, "error_message": "method not allowed"})
	})
	go manager.GC()
	return engine
}

func CORSMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/api/") && authCORS(c) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, accept, origin, Cache-Control, X-Requested-With, sentry-trace, baggage,"+cfg.TokenHeader)
			c.Writer.Header().Set("Access-Control-Expose-Headers", cfg.TokenHeader)
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
			if cfg.Env != config.ApplicationEnvLocal {
				c.Writer.Header().Set("Access-Control-Max-Age", "86400")
			}
			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(204)
				return
			}
		}
		c.Next()
	}
}

func authCORS(c *gin.Context) bool {
	domain := c.Request.Host
	return domain == "localhost"
}
