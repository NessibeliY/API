package document

import (
	"context"
	"database/sql"
	"log"
	"time"

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
	INSERT INTO documents (title, content, image_path, author_id, date_created, date_expired)
	VALUES ($1, $2, $3, $4, $5, $6);`

	_, err := d.db.ExecContext(ctx, query, document.Title, document.Content,
		document.ImagePath, document.AuthorID, document.DateCreated, document.DateExpired)
	if err != nil {
		return err
	}

	return nil
}

func (d *DocumentDatabase) GetDocumentIDByTitle(ctx context.Context, title string) (uint64, error) {
	query := `
	SELECT id
	FROM documents
	WHERE title=$1;`
	var documentID uint64

	err := d.db.QueryRowContext(ctx, query, title).Scan(&documentID)
	return documentID, err
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

func (d *DocumentDatabase) ReadDocument(ctx context.Context, id uint64) (*models.Document, error) {
	query := `
	SELECT title, content, image_path, author_id, date_created, date_expired
	FROM documents
	WHERE id=$1;`

	document := &models.Document{ID: id}

	err := d.db.QueryRowContext(ctx, query, id).Scan(
		&document.Title,
		&document.Content,
		&document.ImagePath,
		&document.AuthorID,
		&document.DateCreated,
		&document.DateExpired,
	)

	return document, err
}

func (d *DocumentDatabase) CheckExpDates(ctx context.Context) ([]*models.Document, error) {
	query := `
	SELECT id, title, date_expired
	FROM documents
	WHERE date_expired BETWEEN NOW() AND NOW() + INTERVAL '7 days';`

	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var documents []*models.Document
	for rows.Next() {
		var (
			id           uint64
			title        string
			date_expired time.Time
		)
		err := rows.Scan(&id, &title, &date_expired)
		if err != nil {
			continue
		}

		log.Println(id, title, date_expired)
		documents = append(documents, &models.Document{ID: id, Title: title, DateExpired: date_expired})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return documents, nil
}
