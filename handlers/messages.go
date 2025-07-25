package handlers

import (
	"chat_app/utils"
	"encoding/hex"
	"errors"
)

func (handler *Handler) storeMsgToDB(chatId, senderId string, payload []byte) error {
	chatObjectId, err := utils.ToObjectId(chatId)
	if err != nil {
		return errors.New(err.Type)
	}

	senderObjectId, err := utils.ToObjectId(senderId)
	if err != nil {
		return errors.New(err.Type)
	}

	ciphered, err2 := handler.Cipher.Encrypt(payload)
	if err2 != nil {
		return err2
	}

	encodedCipher := hex.EncodeToString(ciphered)

	if _, err := handler.Models.Message.Create(chatObjectId, senderObjectId, []byte(encodedCipher)); err != nil {
		return err
	}

	return nil
}
