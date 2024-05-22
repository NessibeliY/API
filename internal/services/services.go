package services

import (
	"time"

	"github.com/NessibeliY/API/internal/database"
	"github.com/NessibeliY/API/internal/dto"
	"github.com/NessibeliY/API/internal/models"
	"github.com/NessibeliY/API/internal/services/document"
	"github.com/NessibeliY/API/internal/services/session"
	"github.com/NessibeliY/API/internal/services/user"
)

type UserServices interface { // TODO SOLID, пустой интерфейс. как интерфейсы го отличаются от других, утиная типизация, контракт (отличие го от С++), под капотом интерфейса есть тип и дата
	SignupUser(*dto.SignupRequest) error
	LoginUser(*dto.LoginRequest) error
}

type SessionServices interface {
	SetSession(string, models.SessionUserClient, time.Duration) error
	GetSession(string, *models.SessionUserClient) error
}

type DocumentServices interface {
	AddInfoAndCreateDocument(*dto.CreateDocumentRequest, time.Time) error
	ShowDocument(*dto.ShowDocumentRequest) (*models.Document, error)
}

type Services struct {
	UserServices     UserServices
	SessionServices  SessionServices
	DocumentServices DocumentServices
}

func NewServices(db *database.Database) *Services {
	return &Services{
		UserServices:     user.NewUserServices(db.UserDatabase),
		SessionServices:  session.NewSessionServices(db.SessionDatabase),
		DocumentServices: document.NewDocumentServices(db.DocumentDatabase),
	}
}
