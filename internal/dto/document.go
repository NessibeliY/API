package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateDocumentRequest struct {
	Title     string `json:"title" validate:"required,min=5"`
	Content   string `json:"content" validate:"required,min=5"`
	ImagePath string `json:"image_path"`
}

type ShowDocumentRequest struct {
	DocumentID int64 `json:"documentID" validate:"required"`
}

type CreateDocumentResponse struct { // TODO response should be sent to client
	BaseResponse
}

type ShowDocumentResponse struct {
	Title       string
	Content     string
	ImagePath   string
	AuthorID    uuid.UUID
	DateCreated time.Time
	BaseResponse
}

type BaseResponse struct {
	Message string
	Status  int
	Err     error
}
