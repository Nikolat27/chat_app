package main

import (
	"chat_app/cipher"
	"chat_app/database"
	"chat_app/database/models"
	"chat_app/handlers"
	"chat_app/webserver"
	"errors"
	"log/slog"

	"github.com/spf13/viper"
)

func main() {
	// load .env
	if err := loadConfig(); err != nil {
		panic(err)
	}

	uri, err := getMongoURI()
	if err != nil {
		panic(err)
	}

	db, err := database.New(uri)
	if err != nil {
		panic(err)
	}

	newModels := models.New(db)

	handlerInstance, err := handlers.New(newModels)
	if err != nil {
		panic(err)
	}

	wsInstance := handlers.WebsocketInit()
	handlerInstance.WebSocket = wsInstance

	cipherInstance := cipher.New()
	handlerInstance.Cipher = cipherInstance

	srv := webserver.New(getPort(), handlerInstance)
	defer srv.Close()

	if err := srv.Run(); err != nil {
		slog.Error("running http webserver", "error", err)
		return
	}
}

func loadConfig() error {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			return errors.New(".env file not found")
		}
		return err
	}

	return nil
}

func getMongoURI() (string, error) {
	uri := viper.GetString("MONGO_URI")
	if uri == "" {
		return "", errors.New("MONGO_URI not set")
	}
	return uri, nil
}

func getPort() string {
	port := viper.GetString("PORT")
	if port == "" {
		return "8000" // default port
	}
	return port
}
