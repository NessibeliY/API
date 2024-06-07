package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/NessibeliY/API/internal/database/document"
	"github.com/NessibeliY/API/internal/database/session"
	"github.com/NessibeliY/API/internal/database/user"
	"github.com/NessibeliY/API/internal/dto"
	"github.com/NessibeliY/API/internal/models"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type UserDatabase interface {
	CreateUser(context.Context, *dto.SignupRequest, []byte) error
	CheckUser(context.Context, *dto.LoginRequest) error
}

type SessionDatabase interface {
	SetSessionData(context.Context, string, models.SessionUserClient, time.Duration) error
	GetSessionData(context.Context, string, *models.SessionUserClient) error
}

type DocumentDatabase interface {
	CreateDocument(context.Context, *models.Document) error
	ReadDocument(context.Context, uint64) (*models.Document, error)
	GetAuthorIDByEmail(context.Context, string) (uuid.UUID, error)
	GetDocumentIDByTitle(context.Context, string) (uint64, error)
}

type Database struct {
	UserDatabase     UserDatabase
	DocumentDatabase DocumentDatabase
	SessionDatabase  SessionDatabase
}

func NewDatabase(db *sql.DB, rdb *redis.Client) *Database {
	return &Database{
		UserDatabase:     user.NewUserDatabase(db),
		DocumentDatabase: document.NewDocumentDatabase(db),
		SessionDatabase:  session.NewSessionDatabase(rdb),
	}
}
