package server

import (
	"fmt"
	"net/http"
)

type Server struct {
	Server *http.Server
	Port   string
}

func New(port string) *Server {
	var srv = &Server{
		Port: port,
	}

	srv.setupServer()

	return srv
}

func (srv *Server) setupServer() {
	router := NewRouter()

	srv.Server = &http.Server{
		Addr:    fmt.Sprintf(":%s", srv.Port),
		Handler: router.CoreRouter,
	}
}

func (srv *Server) Run() error {
	return srv.Server.ListenAndServe()
}

func (srv *Server) RunHttps(certFile, keyFile string) error {
	return srv.Server.ListenAndServeTLS(certFile, keyFile)
}

func (srv *Server) Close() error {
	return srv.Server.Close()
}
