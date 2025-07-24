package main

import (
	"chat_app/database"
	"chat_app/database/models"
	"chat_app/handlers"
	"chat_app/server"
	"errors"
	"log/slog"
	"os"
)

func main() {
	uri, err := getMongoURI()
	if err != nil {
		panic(err)
	}

	db, err := database.New(uri)
	if err != nil {
		panic(err)
	}

	newModels := models.New(db)

	handlerInstance := handlers.New(newModels)

	srv := server.New(getPort(), handlerInstance)
	defer srv.Close()

	if err := srv.Run(); err != nil {
		slog.Error("running http server", "error", err)
		return
	}
}

func getMongoURI() (string, error) {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		return "", errors.New("MONGO_URI env variable does not exist")
	}

	return uri, nil
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		return "8000" // default port
	}

	return port
}
