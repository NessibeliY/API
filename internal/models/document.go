package models

import (
	"time"

	"github.com/google/uuid"
)

type Document struct {
	ID          uint64
	Title       string
	Content     string
	ImagePath   string
	AuthorID    uuid.UUID
	DateCreated time.Time
	DateExpired time.Time
}
