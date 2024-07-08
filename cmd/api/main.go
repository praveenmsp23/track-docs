package main

import (
	"context"
	"github.com/drone/signal"
	"github.com/praveenmsp23/trackdocs/pkg/logger"
)

func main() {
	logger.Init()
	defer logger.Sync()
	ctx := signal.WithContext(
		context.Background(),
	)
	server, e := Initialize(ctx)
	if e != nil {
		logger.Fatal(e)
	}
	logger.Infof("Starting server on %s:%s", server.Listen, server.Port)
	if err := server.ListenAndServe(ctx); err != nil {
		logger.Fatalf("could not run server: %v", err)
	}
}
