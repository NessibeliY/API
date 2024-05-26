package client

import (
	"net/http"
	"time"

	"github.com/NessibeliY/API/internal/dto"
	"github.com/NessibeliY/API/internal/validator"
	"github.com/gin-gonic/gin"
)

func (c *Client) CreateClientDocument(ctx *gin.Context) {
	date := time.Now()
	userEmail, err := ctx.Cookie("user")
	baseResponse := dto.BaseResponse{
		Message: "failed to extract cookie \"user\" from request",
		Status:  http.StatusInternalServerError,
		Err:     err,
	}
	if err != nil {
		ctx.JSON(baseResponse.Status, baseResponse)
		return
	}

	var request dto.CreateDocumentRequest
	err = ctx.ShouldBindJSON(&request)
	baseResponse = dto.BaseResponse{
		Message: "failed to bind JSON",
		Status:  http.StatusBadRequest,
		Err:     err,
	}
	if err != nil {
		ctx.JSON(baseResponse.Status, baseResponse)
		return
	}

	validate := validator.New()
	err = validate.Struct(request)
	baseResponse = dto.BaseResponse{
		Message: "failed to validate request",
		Status:  http.StatusBadRequest,
		Err:     err,
	}
	if err != nil {
		ctx.JSON(baseResponse.Status, baseResponse)
		return
	}

	err = c.documentServices.AddInfoAndCreateDocument(&request, date, userEmail)
	baseResponse = dto.BaseResponse{
		Message: "failed to create document",
		Status:  http.StatusBadRequest,
		Err:     err,
	}
	if err != nil {
		ctx.JSON(baseResponse.Status, baseResponse)
		return
	}

	createDocumentResponse := dto.CreateDocumentResponse{
		BaseResponse: dto.BaseResponse{
			Message: "document successfully created",
			Status:  http.StatusOK,
			Err:     nil,
		},
	}

	ctx.JSON(createDocumentResponse.Status, createDocumentResponse)
}

func (c *Client) ShowClientDocument(ctx *gin.Context) {
	var request dto.ShowDocumentRequest
	err := ctx.ShouldBindJSON(&request)
	baseResponse := dto.BaseResponse{ // TODO think about how to fill baseResponse
		Message: "failed to bind JSON",
		Status:  http.StatusBadRequest,
		Err:     err,
	}
	if err != nil {
		ctx.JSON(baseResponse.Status, baseResponse)
		return
	}

	validate := validator.New()
	err = validate.Struct(request)
	baseResponse = dto.BaseResponse{
		Message: "failed to validate request",
		Status:  http.StatusBadRequest,
		Err:     err,
	}
	if err != nil {
		ctx.JSON(baseResponse.Status, baseResponse)
		return
	}

	document, err := c.documentServices.ShowDocument(&request)
	baseResponse = dto.BaseResponse{
		Message: "failed to show document",
		Status:  http.StatusBadRequest,
		Err:     err,
	}
	if err != nil {
		ctx.JSON(baseResponse.Status, baseResponse)
		return
	}

	showDocumentResponse := dto.ShowDocumentResponse{
		Title:       document.Title,
		Content:     document.Content,
		ImagePath:   document.ImagePath,
		AuthorID:    document.AuthorID,
		DateCreated: document.DateCreated,
		BaseResponse: dto.BaseResponse{
			Message: "document successfully read",
			Status:  http.StatusOK,
			Err:     nil,
		},
	}

	ctx.JSON(showDocumentResponse.Status, showDocumentResponse)
}
