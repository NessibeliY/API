package telegram

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/NessibeliY/API/internal/dto"
	"github.com/go-co-op/gocron"
)

var (
	mu                sync.Mutex
	previousDocuments = make(map[uint64]bool)
)

func (t *TelegramBot) Tgjob(scheduler *gocron.Scheduler) {
	scheduler.Every(5).Seconds().Do(func() {
		expDocuments, err := t.CheckExpDate()
		if err != nil {
			log.Printf("Error checking expiration dates: %v\n", err)
			return
		}

		mu.Lock()
		var newDocs []dto.ExpDocument
		for _, doc := range expDocuments {
			if !previousDocuments[doc.ID] {
				newDocs = append(newDocs, doc)
				previousDocuments[doc.ID] = true
			}
		}
		mu.Unlock()
		// TODO сделать поле в таблице документов, отправлено ли было уведомление
		if len(newDocs) > 0 {
			message, err := json.Marshal(newDocs)
			if err != nil {
				log.Printf("Error marshaling new documents: %v\n", err)
				return
			}

			botToken := "7190440176:AAF4MaUl12HIUPt94D148w7Pi8rXWKkEJMo"
			chatID := "443146012"
			t.SendTelegramNotification(botToken, chatID, string(message))
			if err != nil {
				log.Printf("Error sending Telegram notification: %v\n", err)
			}
		}
	})
}
