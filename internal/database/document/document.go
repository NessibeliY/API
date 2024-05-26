package document

import (
	"context"
	"database/sql"

	"github.com/NessibeliY/API/internal/models"
	"github.com/google/uuid"
)

type DocumentDatabase struct {
	db *sql.DB
}

func NewDocumentDatabase(db *sql.DB) *DocumentDatabase {
	return &DocumentDatabase{
		db: db,
	}
}

func (d *DocumentDatabase) CreateDocument(ctx context.Context, document *models.Document) error {
	query := `
	INSERT INTO documents (id, title, content, image_path, author_id, date_created)
	VALUES ($1, $2, $3, $4, $5, $6);`

	_, err := d.db.ExecContext(ctx, query, document.ID, document.Title,
		document.Content, document.ImagePath, document.AuthorID, document.DateCreated)
	if err != nil {
		return err
	}
	return nil
}

func (d *DocumentDatabase) GetAuthorIDByEmail(ctx context.Context, userEmail string) (uuid.UUID, error) {
	var authorID uuid.UUID

	query := `
	SELECT id
	FROM users
	WHERE email=$1;`

	err := d.db.QueryRowContext(ctx, query, userEmail).Scan(&authorID)
	if err != nil {
		return uuid.Nil, err
	}
	return authorID, nil
}

func (d *DocumentDatabase) ReadDocument(ctx context.Context, title string) (*models.Document, error) {
	query := `
	SELECT id, content, image_path, author_id, date_created
	FROM documents
	WHERE title=$1;`

	document := &models.Document{Title: title}

	err := d.db.QueryRowContext(ctx, query, title).Scan(
		&document.ID,
		&document.Content,
		&document.ImagePath,
		&document.AuthorID,
		&document.DateCreated,
	)

	return document, err
}
