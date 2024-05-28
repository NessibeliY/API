package session

import (
	"context"
	"time"

	"github.com/NessibeliY/API/internal/models"
)

type SessionDatabase interface {
	SetSessionData(context.Context, string, models.SessionUserClient, time.Duration) error
	GetSessionData(context.Context, string, *models.SessionUserClient) error
}

type SessionServices struct {
	sessionDatabase SessionDatabase
}

func NewSessionServices(sessionDatabase SessionDatabase) *SessionServices {
	return &SessionServices{
		sessionDatabase: sessionDatabase,
	}
}

func (ss *SessionServices) SetSession(key string, value models.SessionUserClient, expiration time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := ss.sessionDatabase.SetSessionData(ctx, key, value, expiration)
	if err != nil {
		return err
	}

	return nil
}

func (ss *SessionServices) GetSession(key string, dest *models.SessionUserClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := ss.sessionDatabase.GetSessionData(ctx, key, dest)
	if err != nil {
		return err
	}

	return nil
}
