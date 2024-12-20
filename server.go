package api

import (
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Start(port string, handler *chi.Mux) error {
	s.httpServer = &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
