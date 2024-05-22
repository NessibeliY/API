package dto

type CreateDocumentRequest struct {
	Title     string `json:"title" validate:"required,min=5"`
	Content   string `json:"content" validate:"required,min=5"`
	ImagePath string `json:"image-path"`
}

type ShowDocumentRequest struct {
	Title string `json:"title" validate:"required"`
}
