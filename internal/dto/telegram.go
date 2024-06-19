package dto

import "time"

type ReqBody struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

type ExpDocument struct {
	ID          uint64
	Title       string
	DateExpired time.Time
}

type ExpDocumentResponse struct {
	ExpDocuments []ExpDocument
	BaseResponse
}

// type BaseResponse struct {
// 	Message string
// 	Status  int
// 	Err     error
// }
