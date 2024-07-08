package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/praveenmsp23/trackdocs/pkg/cache"
	"github.com/praveenmsp23/trackdocs/pkg/config"
	"github.com/praveenmsp23/trackdocs/pkg/lock"
	"github.com/praveenmsp23/trackdocs/pkg/models"
	"github.com/praveenmsp23/trackdocs/pkg/service"
	"github.com/praveenmsp23/trackdocs/pkg/store"
	"github.com/praveenmsp23/trackdocs/pkg/token"
)

type Api struct {
	cfg          *config.Config
	tokenManager *token.Manager
	repo         *store.Store
	srv          *service.Service
	redisLock    *lock.RedisLock
	cache        *cache.Cache
}

func (s *Api) Routes(router *gin.RouterGroup) {

	// Public endpoints
	router.Use(Errors(s.cfg))
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, models.NewSuccessResponse("ok"))
	})

	// Account endpoints
	account := router.Group("/account")
	account.Use(AuthMiddleware(s.repo, s.tokenManager))
	account.Use(RateLimitMiddleware(100, s.cache)) // 100 requests per minute per account
	{
		account.GET("/me", HandleGetAccount())
		account.POST("/me/update", HandleAccountUpdate(s.repo))
		account.POST("/logout", HandleAccountLogout(s.tokenManager))
	}
}

func NewApi(cfg *config.Config, store *store.Store, token *token.Manager, lock *lock.RedisLock, srv *service.Service, cache *cache.Cache) (*Api, error) {
	return &Api{cfg: cfg, repo: store, tokenManager: token, redisLock: lock, srv: srv, cache: cache}, nil
}
