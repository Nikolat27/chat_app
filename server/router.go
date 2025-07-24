package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	CoreRouter *chi.Mux
}

func NewRouter() *Router {
	routerInstance := chi.NewRouter()

	routerInstance.Use(middleware.Logger)

	// prefix api
	routerInstance.Route("/api", func(r chi.Router) {

	})

	return &Router{
		CoreRouter: routerInstance,
	}
}
