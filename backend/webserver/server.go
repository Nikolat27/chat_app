package WebServer

import (
	"chat_app/handlers"
	"fmt"
	"net/http"
)

type Server struct {
	Server *http.Server
	Port   string
}

func New(port string, handler *handlers.Handler) *Server {
	var srv = &Server{
		Port: port,
	}

	srv.setupServer(handler)

	return srv
}

func (srv *Server) setupServer(handler *handlers.Handler) {
	router := NewRouter(handler)

	srv.Server = &http.Server{
		Addr:    fmt.Sprintf(":%s", srv.Port),
		Handler: router.CoreRouter,
	}
}

func (srv *Server) Run() error {
	fmt.Println("application started")
	return srv.Server.ListenAndServe()
}

func (srv *Server) RunHttps(certFile, keyFile string) error {
	return srv.Server.ListenAndServeTLS(certFile, keyFile)
}

func (srv *Server) Close() error {
	return srv.Server.Close()
}
