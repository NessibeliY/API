package telegram

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/NessibeliY/API/internal/dto"
)

func SendTelegramNotification(botToken, chatID, message string) error {
	log.Println("Preparing to send Telegram notification")

	apiUrl := "https://api.telegram.org/bot" + botToken + "/sendMessage"

	// Create request
	reqBodyValue := dto.ReqBody{
		ChatID: chatID,
		Text:   message,
	}

	reqBodyValueJson, err := json.Marshal(reqBodyValue)
	if err != nil {
		return errors.New("Failed to marshal request body for telegram: " + err.Error())
	}

	reqBodyJson := bytes.NewReader(reqBodyValueJson)
	if err != nil {
		return errors.New("Failed to make reader for request body for telegram: " + err.Error())
	}

	req, err := http.NewRequest(http.MethodPost, apiUrl, reqBodyJson)
	if err != nil {
		return errors.New("Failed to create new request for telegram: " + err.Error())
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	// make request
	res, err := client.Do(req)
	if err != nil {
		return errors.New("Failed to send telegram notification: " + err.Error())
	}

	if res.StatusCode/100 != 2 {
		log.Println("Failed to send Telegram notification", res.Status)
		return errors.New("Failed to send Telegram notification")
	}

	log.Println("Notification sent to Telegram")
	return nil
}
