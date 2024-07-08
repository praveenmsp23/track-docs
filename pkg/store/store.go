package store

import (
	"github.com/praveenmsp23/trackdocs/pkg/cache"
	"github.com/praveenmsp23/trackdocs/pkg/config"
	"gorm.io/gorm"
)

const SearchLimit = 10

// Store one stop for stores
type Store struct {
	AccountStore *accountStore
}

// NewStore create all the stores
func NewStore(conn *gorm.DB, cache *cache.Cache, cfg *config.Config) (*Store, error) {
	repo := &Store{
		AccountStore: newAccountStore(conn, cache, cfg),
	}
	repo.AccountStore.repo = repo
	return repo, nil
}
