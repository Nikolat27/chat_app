package HttpServer

import (
	"chat_app/handlers"
	"fmt"
	"net/http"
)

type HttpServer struct {
	Server *http.Server
	Port   string
}

func New(port string, handler *handlers.Handler) *HttpServer {
	var srv = &HttpServer{
		Port: port,
	}

	srv.setupServer(handler)

	return srv
}

func (srv *HttpServer) setupServer(handler *handlers.Handler) {
	router := NewRouter(handler)

	srv.Server = &http.Server{
		Addr:    fmt.Sprintf(":%s", srv.Port),
		Handler: router.CoreRouter,
	}
}

func (srv *HttpServer) Run() error {
	fmt.Println("application started")
	return srv.Server.ListenAndServe()
}

func (srv *HttpServer) RunHttps(certFile, keyFile string) error {
	return srv.Server.ListenAndServeTLS(certFile, keyFile)
}

func (srv *HttpServer) Close() error {
	return srv.Server.Close()
}
