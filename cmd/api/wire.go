//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"github.com/google/wire"
	"github.com/praveenmsp23/trackdocs/pkg/cache"
	"github.com/praveenmsp23/trackdocs/pkg/config"
	"github.com/praveenmsp23/trackdocs/pkg/db"
	"github.com/praveenmsp23/trackdocs/pkg/lock"
	"github.com/praveenmsp23/trackdocs/pkg/server"
	"github.com/praveenmsp23/trackdocs/pkg/service"
	"github.com/praveenmsp23/trackdocs/pkg/store"
	"github.com/praveenmsp23/trackdocs/pkg/token"
)

func Initialize(ctx context.Context) (*server.Server, error) {
	wire.Build(
		config.NewConfig,
		cache.NewCache,
		lock.NewRedisLock,
		db.NewDB,
		service.NewService,
		store.NewStore,
		token.NewManager,
		serverSet,
		server.InitServer,
	)
	return &server.Server{}, nil
}
