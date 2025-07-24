package handlers

import (
	"chat_app/database/models"
	"chat_app/paseto"
)

type Handler struct {
	Models *models.Models
	Paseto *paseto.Maker
}

func New(models *models.Models) (*Handler, error) {
	pasetoInstance, err := paseto.New()
	if err != nil {
		return nil, err
	}

	var handler = &Handler{
		Models: models,
		Paseto: pasetoInstance,
	}

	return handler, nil
}
