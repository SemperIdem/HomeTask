package server

import (
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
)

type Server struct {
	router *http.ServeMux
	srv    *http.Server

	ctrl   NumbersFetcher
	result chan []int
}

type NumbersFetcher interface {
	ProcessUrls(urls []string)
	Receive() chan []int
}

//go:generate mockgen -destination=./mocks/mocks_numbersFetcher.go --build_flags=--mod=mod -package=mocks homeTask/server NumbersFetcher

func New(ctx context.Context, ctrl NumbersFetcher, result chan []int) (*Server, func()) {
	s := &Server{
		router: http.NewServeMux(),
		result: result,
		ctrl:   ctrl,
	}

	return s, func() {
		s.Shutdown(ctx)
	}
}

func (srv *Server) ListenAndServeHTTP(addr string) error {
	log.Printf("Starting http listener on %s", addr)

	srv.srv = &http.Server{
		Addr:    addr,
		Handler: srv.router,
	}

	ln, err := net.Listen("tcp", srv.srv.Addr)
	if err != nil {
		return fmt.Errorf("listen error: %w", err)
	}

	go func() {
		if err := srv.srv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("ListenAndServeHTTP: srv.Serve", err)
		}
	}()

	return nil
}

func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.router.ServeHTTP(w, r)
}

func (srv *Server) Shutdown(ctx context.Context) {
	if err := srv.srv.Shutdown(ctx); err != nil {
		log.Error("http server shutdown with error: ", err)
	}
}

func (srv *Server) SetUpRoutes() {
	srv.router.HandleFunc("/numbers", Numbers(srv.ctrl))

}

func (srv *Server) GetRouter() *http.ServeMux {
	return srv.router
}
