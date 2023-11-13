package server

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/d-jackalope/L0/internal/app"
)

type Server struct {
	App      *app.Application
	ErrChan  chan error
	DoneChan chan struct{}
}

func New(app *app.Application) *Server {
	server := &Server{
		App:      app,
		ErrChan:  make(chan error),
		DoneChan: make(chan struct{}),
	}
	server.App.Srv.Handler = server.routers()
	return server
}

func (s *Server) Run(ctx context.Context) {
	go func() {
		if err := s.App.Srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.ErrChan <- err
		}
	}()

	s.App.Log.Info("Starting the server on %s\n", s.App.Srv.Addr)
	<-ctx.Done()
	s.close()
	close(s.DoneChan)
}

func (s *Server) close() {
	ctx, cansel := context.WithTimeout(context.Background(), time.Second*5)
	defer cansel()
	if err := s.App.Srv.Shutdown(ctx); err != nil {
		trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
		s.App.Log.Output(2, trace)
	} else {
		s.App.Log.Info("The server was closed")
	}
}
