package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/praveenmsp23/trackdocs/pkg/config"
	"github.com/praveenmsp23/trackdocs/pkg/models"
	"github.com/praveenmsp23/trackdocs/pkg/store"
	"github.com/praveenmsp23/trackdocs/pkg/token"
)

type Health struct {
	Config *config.Config
	Store  *store.Store
	Token  *token.Manager
}

func (s *Health) Routes(router *gin.RouterGroup) {
	router.GET("/", func(c *gin.Context) {
		build := gin.H{
			"commit": config.BuildSHA,
			"branch": config.BuildBranch,
			"time":   config.BuildTime,
		}
		c.JSON(http.StatusOK, models.NewSuccessResponse(gin.H{"status": "ok", "build": build}))
	})
}

func NewHealth(cfg *config.Config, store *store.Store, token *token.Manager) (*Health, error) {
	return &Health{Config: cfg, Store: store, Token: token}, nil
}
