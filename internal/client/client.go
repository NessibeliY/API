package client

import (
	"time"

	"github.com/NessibeliY/API/internal/dto"
	"github.com/NessibeliY/API/internal/models"
	"github.com/NessibeliY/API/internal/services"
)

type DocumentServices interface {
	AddInfoAndCreateDocument(*dto.CreateDocumentRequest, time.Time, string) error
	ShowDocument(*dto.ShowDocumentRequest) (*models.Document, error)
}

type UserServices interface { // TODO SOLID, пустой интерфейс. как интерфейсы го отличаются от других, утиная типизация, контракт (отличие го от С++), под капотом интерфейса есть тип и дата
	SignupUser(*dto.SignupRequest) error
	LoginUser(*dto.LoginRequest) error
}

type SessionServices interface {
	SetSession(string, models.SessionUserClient, time.Duration) error
	GetSession(string, *models.SessionUserClient) error
}

type Client struct {
	userServices     UserServices
	sessionServices  SessionServices
	documentServices DocumentServices
}

func NewClient(srv *services.Services) *Client {
	return &Client{
		userServices:     srv.UserServices,
		sessionServices:  srv.SessionServices,
		documentServices: srv.DocumentServices,
	}
}
