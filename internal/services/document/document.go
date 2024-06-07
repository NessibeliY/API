package document

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/NessibeliY/API/internal/dto"
	"github.com/NessibeliY/API/internal/models"
	"github.com/google/uuid"
)

type DocumentDatabase interface {
	CreateDocument(context.Context, *models.Document) error
	ReadDocument(context.Context, uint64) (*models.Document, error)
	GetAuthorIDByEmail(context.Context, string) (uuid.UUID, error)
	GetDocumentIDByTitle(context.Context, string) (uint64, error)
}

type DocumentServices struct {
	documentDatabase DocumentDatabase
}

func NewDocumentServices(documentDatabase DocumentDatabase) *DocumentServices {
	return &DocumentServices{
		documentDatabase: documentDatabase,
	}
}

func (ds *DocumentServices) CreateDocument(request *dto.CreateDocumentRequest, date time.Time, userEmail string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // TODO move 3*time.Second to value/constants
	defer cancel()
	fmt.Println("\nin services\n", request)

	// Check if document already exists
	_, err := ds.documentDatabase.GetDocumentIDByTitle(ctx, request.Title)
	if err != sql.ErrNoRows {
		return errors.New("Such title already exists")
	}
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	authorID, err := ds.documentDatabase.GetAuthorIDByEmail(ctx, userEmail)
	if err == sql.ErrNoRows {
		return errors.New("The user is not logged in")
	}
	if err != nil {
		return err
	}

	processedRequest := &models.Document{
		Title:       request.Title,
		Content:     request.Content,
		ImagePath:   request.ImagePath,
		DateCreated: date,
		DateExpired: request.ExpirationDate,
		AuthorID:    authorID,
	}

	// Create document
	err = ds.documentDatabase.CreateDocument(ctx, processedRequest)
	if err != nil {
		return err
	}

	// Identify author, size, time of creation

	return nil
}

func (ds *DocumentServices) GetDocument(request *dto.ShowDocumentRequest) (*models.Document, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	document, err := ds.documentDatabase.ReadDocument(ctx, request.DocumentID)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("Such title does not exist")
	}

	if err != nil {
		return nil, err
	}

	return document, nil
}
