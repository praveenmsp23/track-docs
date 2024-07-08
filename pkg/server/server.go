package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/praveenmsp23/trackdocs/pkg/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/praveenmsp23/trackdocs/pkg/config"
	"golang.org/x/sync/errgroup"
)

// A Server defines parameters for running an HTTP server.
type Server struct {
	Port    string
	Listen  string
	Handler http.Handler
}

const timeoutGracefulShutdown = 5 * time.Second

// ListenAndServe initializes a server to respond to HTTP network requests.
func (s Server) ListenAndServe(ctx context.Context) error {
	err := s.listenAndServe(ctx)
	if errors.Is(err, http.ErrServerClosed) {
		err = nil
	}
	return err
}

func (s Server) listenAndServe(ctx context.Context) error {
	var g errgroup.Group

	s1 := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", s.Listen, s.Port),
		Handler: s.Handler,
	}
	g.Go(func() error {
		<-ctx.Done()

		ctxShutdown, cancelFunc := context.WithTimeout(context.Background(), timeoutGracefulShutdown)
		defer cancelFunc()
		return s1.Shutdown(ctxShutdown)
	})
	g.Go(s1.ListenAndServe)
	return g.Wait()
}

func InitServer(cfg *config.Config, engine *gin.Engine) (*Server, error) {
	return &Server{Port: cfg.Port, Listen: cfg.Listen, Handler: engine}, nil
}

func HandleFunc(handler func(*models.TrackDocsContext)) func(*gin.Context) {
	return func(c *gin.Context) {
		handler(models.NewTrackDocsContext(c))
	}
}
