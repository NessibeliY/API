package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/NessibeliY/API/internal/dto"
	"github.com/NessibeliY/API/internal/models"
	"github.com/google/uuid"
)

// TODO create chatid table that sends info to specific user, tdid and userid as foreignt key, use JOIN
// TODO 4 links save and learn, goroutines
// TODO если юзер не ответил кодом, то отправить ему еще уведомление
// TODO read about migrations
type DocumentDatabase interface {
	CreateDocument(context.Context, *models.Document) error
	ReadDocument(context.Context, uint64) (*models.Document, error)
	GetAuthorIDByEmail(context.Context, string) (uuid.UUID, error)
	GetDocumentIDByTitle(context.Context, string) (uint64, error)
	CheckExpDates(context.Context) ([]*models.Document, error)
}

type TelegramBot struct {
	documentDatabase DocumentDatabase
}

func NewTelegramBot(documentDatabase DocumentDatabase) *TelegramBot {
	return &TelegramBot{
		documentDatabase: documentDatabase,
	}
}

func (t *TelegramBot) CheckExpDate() ([]dto.ExpDocument, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	documents, err := t.documentDatabase.CheckExpDates(ctx)
	if err != nil {
		return nil, errors.New("Failed to retrieve data from db: " + err.Error())
	}

	expDocuments := make([]dto.ExpDocument, 0)

	for _, document := range documents {
		expDocument := dto.ExpDocument{
			ID:          document.ID,
			Title:       document.Title,
			DateExpired: document.DateExpired,
		}
		expDocuments = append(expDocuments, expDocument)
	}
	// if len(expDocuments) == 0 {
	// 	return nil, nil
	// }

	// expDocumentResponse := dto.ExpDocumentResponse{
	// 	ExpDocuments: expDocuments,
	// 	BaseResponse: baseResponse,
	// }

	// responseJSON, err := json.Marshal(expDocumentResponse)
	// if err != nil {
	// 	return nil, errors.New("failed to marshal response for telegram" + err.Error())
	// }

	return expDocuments, nil
}

func (t *TelegramBot) SendTelegramNotification(botToken, chatID, message string) error {
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
