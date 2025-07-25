package server

import (
	"chat_app/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	CoreRouter *chi.Mux
}

func NewRouter(handler *handlers.Handler) *Router {
	routerInstance := chi.NewRouter()

	routerInstance.Use(middleware.Logger)

	// prefix api
	routerInstance.Route("/api", func(r chi.Router) {
		r.Post("/auth/register", handler.Register)
		r.Post("/auth/login", handler.Login)

		r.Get("/websocket/chat/add", handler.AddChatWebsocket)
	})

	return &Router{
		CoreRouter: routerInstance,
	}
}
