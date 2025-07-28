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

	routerInstance.Use(CheckCorsOrigin)
	routerInstance.Use(middleware.Logger)

	// prefix api
	routerInstance.Route("/api", func(r chi.Router) {

		r.Post("/register", handler.Register)
		r.Post("/login", handler.Login)

		r.Get("/user/search", handler.SearchUser)
		r.Delete("/user/delete", handler.DeleteUser)
		r.Post("/user/upload-avatar", handler.UploadAvatar)
		r.Get("/user/get-chats", handler.GetUserChats)
		r.Get("/user/get-secret-chats", handler.GetUserSecretChats)

		r.Post("/chat/create", handler.CreateChat)
		r.Post("/chat/upload/{chat_id}/{receiver_id}", handler.UploadChatImage)
		r.Get("/chat/get/{chat_id}/messages", handler.GetChatMessages)
		r.Delete("/chat/delete/{chat_id}", handler.DeleteChat)
		// chat websocket
		r.Get("/websocket/chat/add/{chat_id}", handler.AddChatWebsocket)

		r.Put("/message/update/{message_id}", handler.EditMessage)
		r.Delete("/message/delete/sender/{message_id}", handler.DeleteMessageForSender)
		r.Delete("/message/delete/all/{message_id}", handler.DeleteMessageForAll)

		r.Post("/group/create", handler.CreateGroup)
		r.Get("/group/join/{invite_link}", handler.JoinGroup)
		r.Get("/group/get/{group_id}", handler.GetGroupMessages)
		r.Delete("/group/remove-user/{group_id}/{user_id}", handler.RemoveUserFromGroup)
		r.Delete("/group/delete/{group_id}", handler.DeleteGroup)
		// group websocket
		r.Get("/websocket/group/add/{group_id}", handler.AddGroupWebsocket)

		r.Post("/save-message/create", handler.CreateSaveMessage)
		r.Get("/save-message/get", handler.GetSaveMessages)
		r.Put("/save-message/update/{message_id}", handler.EditSaveMessage)
		r.Delete("/save-message/delete", handler.DeleteSaveMessage)

		r.Post("/secret-chat/create", handler.CreateSecretChat)
		r.Get("/secret-chat/get/{secret_chat_id}/messages", handler.GetSecretChatMessages)
		r.Delete("/secret-chat/delete/{secret_chat_id}", handler.DeleteSecretChat)
		r.Post("/secret-chat/add-symmetric-key/{secret_chat_id}", handler.UpdateSecretChat)
		r.Post("/secret-chat/approve/{secret_chat_id}", handler.ApproveSecretChat)
		// chat websocket
		r.Get("/websocket/secret-chat/add/{secret_chat_id}", handler.AddSecretChatWebsocket)

	})

	fs := http.FileServer(http.Dir("./uploads"))
	routerInstance.Handle("/static/*", http.StripPrefix("/static/", fs))

	return &Router{
		CoreRouter: routerInstance,
	}
}
