package service

import (
	"github.com/praveenmsp23/trackdocs/pkg/config"
	"github.com/praveenmsp23/trackdocs/pkg/store"
)

// Service one stop for all the services
type Service struct {
}

// NewService create all the services
func NewService(cfg *config.Config, repo *store.Store) (*Service, error) {
	return &Service{}, nil
}
