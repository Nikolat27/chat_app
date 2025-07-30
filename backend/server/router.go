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

		// Authentication
		r.Post("/register", handler.Register)
		r.Post("/login", handler.Login)

		// Members
		r.Get("/user/search", handler.SearchUser)
		r.Get("/user/get/{user_id}", handler.GetUser)
		r.Delete("/user/delete", handler.DeleteUser)
		r.Post("/user/upload-avatar", handler.UploadAvatar)
		r.Get("/user/get-chats", handler.GetUserChats)
		r.Get("/user/get-secret-chats", handler.GetUserSecretChats)
		r.Get("/user/get-groups", handler.GetUserGroups)
		r.Get("/user/get-secret-groups", handler.GetUserSecretGroups)

		// Chats
		r.Post("/chat/create", handler.CreateChat)
		r.Post("/chat/upload/{chat_id}/{receiver_id}", handler.UploadChatImage)
		r.Get("/chat/get/{chat_id}/messages", handler.GetChatMessages)
		r.Delete("/chat/delete/{chat_id}", handler.DeleteChat)
		// chat websocket
		r.Get("/websocket/chat/add/{chat_id}", handler.AddChatWebsocket)

		// Messages
		r.Put("/message/update/{message_id}", handler.EditMessage)
		r.Delete("/message/delete/sender/{message_id}", handler.DeleteMessageForSender)
		r.Delete("/message/delete/all/{message_id}", handler.DeleteMessageForAll)

		// Groups
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
		// group websocket
		r.Get("/websocket/group/add/{group_id}", handler.AddGroupWebsocket)

		// Secret Groups
		r.Post("/secret-group/create", handler.CreateSecretGroup)
		r.Put("/secret-group/update/{secret_group_id}", handler.UpdateSecretGroup)
		r.Get("/secret-group/join/{invite_link}", handler.JoinSecretGroup)
		r.Post("/secret-group/ban/{secret_group_id}", handler.BanMemberFromSecretGroup)
		r.Post("/secret-group/unban/{secret_group_id}", handler.UnBanMemberFromSecretGroup)
		r.Get("/secret-group/get/{secret_group_id}/messages", handler.GetSecretGroupMessages)
		r.Get("/secret-group/get/{secret_group_id}/members", handler.GetSecretGroupMembers)
		r.Delete("/secret-group/leave/{secret_group_id}", handler.LeaveSecretGroup)
		r.Delete("/secret-group/remove-user/{secret_group_id}/{user_id}", handler.RemoveUserFromSecretGroup)
		r.Delete("/secret-group/delete/{secret_group_id}", handler.DeleteSecretGroup)
		// secret group websocket
		r.Get("/websocket/secret-group/add/{secret_group_id}", handler.AddSecretGroupWebsocket)

		// Save Messages
		r.Post("/save-message/create", handler.CreateSaveMessage)
		r.Get("/save-message/get", handler.GetSaveMessages)
		r.Put("/save-message/update/{message_id}", handler.EditSaveMessage)
		r.Delete("/save-message/delete/{message_id}", handler.DeleteSaveMessage)

		// Secret Chats
		r.Get("/secret-chat/get/{secret_chat_id}", handler.GetSecretChat)
		r.Post("/secret-chat/create", handler.CreateSecretChat)
		r.Get("/secret-chat/get/{secret_chat_id}/messages", handler.GetSecretChatMessages)
		r.Delete("/secret-chat/delete/{secret_chat_id}", handler.DeleteSecretChat)
		r.Post("/secret-chat/add-public-key/{secret_chat_id}", handler.UploadSecretChatPublicKey)
		r.Post("/secret-chat/add-symmetric-key/{secret_chat_id}", handler.UploadSecretChatSymmetricKey)
		r.Post("/secret-chat/approve/{secret_chat_id}", handler.ApproveSecretChat)
		// chat websocket
		r.Get("/websocket/secret-chat/add/{secret_chat_id}", handler.AddSecretChatWebsocket)

		// Approvals
		r.Get("/received-approvals/get/", handler.GetReceivedApprovals)
		r.Get("/sent-approvals/get/", handler.GetSentApprovals)
		r.Post("/approvals/submit/{invite_link}", handler.CreateApproval)
		r.Put("/approvals/edit-status/{approval_id}", handler.EditApprovalStatus)
		r.Delete("/approvals/delete/{approval_id}", handler.DeleteApproval)

	})

	fs := http.FileServer(http.Dir("./uploads"))
	routerInstance.Handle("/static/*", http.StripPrefix("/static/", fs))

	return &Router{
		CoreRouter: routerInstance,
	}
}
