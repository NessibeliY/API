package document

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/NessibeliY/API/internal/dto"
	"github.com/NessibeliY/API/internal/models"
	"github.com/google/uuid"
)

type DocumentDatabase interface {
	CreateDocument(context.Context, *models.Document) error
	ReadDocument(context.Context, string) (*models.Document, error)
}

type DocumentServices struct {
	documentDatabase DocumentDatabase
}

func NewDocumentServices(documentDatabase DocumentDatabase) *DocumentServices {
	return &DocumentServices{
		documentDatabase: documentDatabase,
	}
}

func (ds *DocumentServices) AddInfoAndCreateDocument(request *dto.CreateDocumentRequest, date time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Check if document already exists
	_, err := ds.documentDatabase.ReadDocument(ctx, request.Title)
	if err != sql.ErrNoRows {
		return fmt.Errorf("Such title already exists")
	}

	processedRequest := &models.Document{
		ID:          uuid.New(),
		Title:       request.Title,
		Content:     request.Content,
		ImagePath:   request.ImagePath,
		DateCreated: date,
	}

	// Create document
	err = ds.documentDatabase.CreateDocument(ctx, processedRequest)
	if err != nil {
		return err
	}

	// Identify author, size, time of creation

	return nil
}

func (ds *DocumentServices) ShowDocument(request *dto.ShowDocumentRequest) (*models.Document, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	document, err := ds.documentDatabase.ReadDocument(ctx, request.Title)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("Such title does not exist")
	}
	if err != nil {
		return nil, err
	}

	return document, nil
}
