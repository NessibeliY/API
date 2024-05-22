package document

import (
	"context"
	"database/sql"

	"github.com/NessibeliY/API/internal/models"
)

type DocumentDatabase struct {
	db *sql.DB
}

func NewDocumentDatabase(db *sql.DB) *DocumentDatabase {
	return &DocumentDatabase{
		db: db,
	}
}

func (ddb *DocumentDatabase) CreateDocument(ctx context.Context, document *models.Document) error {
	query := `
	INSERT INTO documents (id, title, content, image_path, author_id, date_created)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id;`

	_, err := ddb.db.ExecContext(ctx, query, document.ID, document.Title,
		document.Content, document.ImagePath, document.AuthorID, document.DateCreated)
	if err != nil {
		return err
	}
	return nil
}

func (ddb *DocumentDatabase) ReadDocument(ctx context.Context, title string) (*models.Document, error) {
	query := `
	SELECT id, content, image_path, author_id, date_created
	FROM documents
	WHERE title=$1;`

	document := &models.Document{Title: title}

	err := ddb.db.QueryRowContext(ctx, query, title).Scan(
		&document.ID,
		&document.Content,
		&document.ImagePath,
		&document.AuthorID,
		&document.DateCreated,
	)

	return document, err
}
