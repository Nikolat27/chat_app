package server

import (
	"chat_app/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
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

		r.Delete("/user/delete", handler.DeleteUser)
		r.Put("/user/upload-avatar", handler.UploadAvatar)

		r.Post("/chat/create", handler.CreateChat)
		r.Post("/chat/upload/{chat_id}/{receiver_id}", handler.UploadChatImage)
		r.Get("/chat/get/{chat_id}", handler.GetChatMessages)
		r.Delete("/chat/delete/{chat_id}", handler.DeleteChat)

		r.Get("/websocket/chat/add/{chat_id}", handler.AddChatWebsocket)

		r.Put("/message/update/{message_id}", handler.EditMessage)
		r.Delete("/message/delete/sender/{message_id}", handler.DeleteMessageForSender)
		r.Delete("/message/delete/receiver/{message_id}", handler.DeleteMessageForReceiver)
		r.Delete("/message/delete/all/{message_id}", handler.DeleteMessageForAll)

		r.Post("/group/create", handler.CreateGroup)
		r.Get("/group/join/{invite_link}", handler.JoinGroup)
		r.Delete("/group/remove-user/{group_id}/{user_id}", handler.RemoveUserFromGroup)
		r.Delete("/group/delete/{group_id}", handler.DeleteGroup)
		
	})

	fs := http.FileServer(http.Dir("./uploads"))
	routerInstance.Handle("/static/*", http.StripPrefix("/static/", fs))

	return &Router{
		CoreRouter: routerInstance,
	}
}
