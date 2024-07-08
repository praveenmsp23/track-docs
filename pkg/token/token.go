package token

import (
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/praveenmsp23/trackdocs/pkg/config"
	"github.com/praveenmsp23/trackdocs/pkg/crypto"
	"github.com/praveenmsp23/trackdocs/pkg/token/base"
	"github.com/praveenmsp23/trackdocs/pkg/token/providers/memory"
	"github.com/praveenmsp23/trackdocs/pkg/token/providers/redis"
)

type Manager struct {
	headerName string     //private header name
	lock       sync.Mutex // protects token
	provider   base.Provider
	lifetime   int64
	cfg        *config.Config
}

func NewManager(cfg *config.Config) (*Manager, error) {
	var provider base.Provider
	if cfg.TokenProvider == "memory" {
		p, err := memory.GetProvider(cfg)
		if err != nil {
			return nil, err
		}
		provider = p
	} else if cfg.TokenProvider == "redis" {
		p, err := redis.GetProvider(cfg)
		if err != nil {
			return nil, err
		}
		provider = p
	} else {
		return nil, fmt.Errorf("token: unknown provide %q (forgotten import?)", cfg.TokenProvider)
	}
	return &Manager{provider: provider, cfg: cfg, headerName: cfg.TokenHeader, lifetime: cfg.TokenLifeTime}, nil
}

// TokenGet get token
func (manager *Manager) TokenGet(c *gin.Context) (token base.Token) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	header := c.Request.Header.Get(manager.headerName)
	if header == "" {
		return nil
	} else {
		tid, _ := url.QueryUnescape(header)
		token, _ = manager.provider.TokenRead(tid)
	}
	return
}

func (manager *Manager) TokenInit(c *gin.Context) (token base.Token) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	tid := manager.tokenId()
	token, _ = manager.provider.TokenInit(tid)
	c.Header(manager.headerName, tid)
	return
}

// TokenDestroy destroy token id
func (manager *Manager) TokenDestroy(c *gin.Context) {
	header := c.Request.Header.Get(manager.headerName)
	if header == "" {
		return
	} else {
		manager.lock.Lock()
		defer manager.lock.Unlock()
		tid, _ := url.QueryUnescape(header)
		manager.provider.TokenDestroy(tid)
	}
}

func (manager *Manager) GC() {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	manager.provider.TokenGC(manager.lifetime)
	time.AfterFunc(time.Duration(manager.lifetime)*time.Second, func() { manager.GC() })
}

func (manager *Manager) tokenId() string {
	return crypto.GenerateId("tok", 32)
}
