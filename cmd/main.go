package main

import (
	"chat_app/handlers"
	"chat_app/server"
	"log/slog"
	"os"
)

func main() {
	port := getPort()

	handlerInstance := handlers.New()

	srv := server.New(port, handlerInstance)
	defer srv.Close()

	if err := srv.Run(); err != nil {
		slog.Error("running http server", "error", err)
		return
	}
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		return "8000" // default port
	}

	return port
}
