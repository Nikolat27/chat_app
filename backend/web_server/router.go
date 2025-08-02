package web_server

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

	routerInstance.Route("/api", func(r chi.Router) {
		getAuthRoutes(r, handler)
		getUserRoutes(r, handler)
		getChatRoutes(r, handler)
		getMessageRoutes(r, handler)
		getGroupRoutes(r, handler)
		getSaveMessageRoutes(r, handler)
		getSecretChatRoutes(r, handler)
		getApprovalRoutes(r, handler)
	})

	fs := http.FileServer(http.Dir("./uploads"))
	routerInstance.Handle("/static/*", http.StripPrefix("/static/", fs))

	return &Router{CoreRouter: routerInstance}
}

func getAuthRoutes(r chi.Router, handler *handlers.Handler) {
	r.Post("/register", handler.Register)
	r.Post("/login", handler.Login)
	r.Get("/logout", handler.Logout)
	r.Get("/auth-check", handler.AuthCheck)
}

func getUserRoutes(r chi.Router, handler *handlers.Handler) {
	r.Get("/user/search", handler.SearchUser)
	r.Get("/user/get/{user_id}", handler.GetUser)
	r.Delete("/user/delete", handler.DeleteUser)
	r.Post("/user/upload-avatar", handler.UploadAvatar)
	r.Get("/user/get-chats", handler.GetUserChats)
	r.Get("/user/get-secret-chats", handler.GetUserSecretChats)
	r.Get("/user/get-groups", handler.GetUserGroups)
}

func getChatRoutes(r chi.Router, handler *handlers.Handler) {
	r.Post("/chat/create", handler.CreateChat)
	r.Post("/chat/upload/{chat_id}/{receiver_id}", handler.UploadChatImage)
	r.Get("/chat/get/{chat_id}/messages", handler.GetChatMessages)
	r.Delete("/chat/delete/{chat_id}", handler.DeleteChat)
	r.Get("/websocket/chat/add/{chat_id}", handler.AddChatWebsocket)
}

func getMessageRoutes(r chi.Router, handler *handlers.Handler) {
	r.Put("/message/update/{message_id}", handler.EditMessage)
	r.Post("/message/upload-chat-image/{chat_id}", handler.UploadImageChatMessage)
	r.Post("/message/upload-group-image/{group_id}", handler.UploadImageGroupMessage)
	r.Delete("/message/delete/sender/{message_id}", handler.DeleteMessageForSender)
	r.Delete("/message/delete/all/{message_id}", handler.DeleteMessageForAll)
}

func getGroupRoutes(r chi.Router, handler *handlers.Handler) {
	r.Post("/group/create", handler.CreateGroup)
	r.Put("/group/update/{group_id}", handler.UpdateGroup)
	r.Get("/group/join/{invite_link}", handler.JoinGroup)
	r.Post("/group/ban/{group_id}", handler.BanMemberFromGroup)
	r.Post("/group/unban/{group_id}", handler.UnBanMemberFromGroup)
	r.Get("/group/get/{group_id}/messages", handler.GetGroupMessages)
	r.Get("/group/get/{group_id}/members", handler.GetGroupMembers)
	r.Delete("/group/leave/{group_id}", handler.LeaveGroup)
	r.Delete("/group/remove-user/{group_id}/{user_id}", handler.RemoveUserFromGroup)
	r.Delete("/group/delete/{group_id}", handler.DeleteGroup)
	r.Get("/websocket/group/add/{group_id}", handler.AddGroupWebsocket)
}

func getSaveMessageRoutes(r chi.Router, handler *handlers.Handler) {
	r.Post("/save-message/create", handler.CreateSaveMessage)
	r.Get("/save-message/get", handler.GetSaveMessages)
	r.Put("/save-message/update/{message_id}", handler.EditSaveMessage)
	r.Delete("/save-message/delete/{message_id}", handler.DeleteSaveMessage)
}

func getSecretChatRoutes(r chi.Router, handler *handlers.Handler) {
	r.Get("/secret-chat/get/{secret_chat_id}", handler.GetSecretChat)
	r.Post("/secret-chat/create", handler.CreateSecretChat)
	r.Get("/secret-chat/get/{secret_chat_id}/messages", handler.GetSecretChatMessages)
	r.Delete("/secret-chat/delete/{secret_chat_id}", handler.DeleteSecretChat)
	r.Post("/secret-chat/add-public-key/{secret_chat_id}", handler.UploadSecretChatPublicKey)
	r.Post("/secret-chat/add-symmetric-key/{secret_chat_id}", handler.UploadSecretChatSymmetricKey)
	r.Post("/secret-chat/approve/{secret_chat_id}", handler.ApproveSecretChat)
	r.Get("/websocket/secret-chat/add/{secret_chat_id}", handler.AddSecretChatWebsocket)
}

func getApprovalRoutes(r chi.Router, handler *handlers.Handler) {
	r.Get("/received-approvals/get/", handler.GetReceivedApprovals)
	r.Get("/sent-approvals/get/", handler.GetSentApprovals)
	r.Post("/approvals/submit/{invite_link}", handler.CreateApproval)
	r.Put("/approvals/edit-status/{approval_id}", handler.EditApprovalStatus)
	r.Delete("/approvals/delete/{approval_id}", handler.DeleteApproval)
}
