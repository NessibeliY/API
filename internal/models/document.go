package models

import (
	"time"

	"github.com/google/uuid"
)

type Document struct {
	ID          uuid.UUID
	Title       string
	Content     string
	ImagePath   string
	AuthorID    uuid.UUID
	DateCreated time.Time
}
