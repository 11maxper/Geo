package httpserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"project/config"
)

type Server struct {
	server *http.Server
}

func NewServer(handler http.Handler, c config.HTTP) *Server {
	return &Server{
		server: &http.Server{
			Addr:         fmt.Sprintf(":%s", c.Port),
			Handler:      handler,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}
}

func (s *Server) Serve() {
	go func() {
		log.Printf("Server is running on port %s", s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()
}

func (s *Server) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.server.Shutdown(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
