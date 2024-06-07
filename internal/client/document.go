package client

import (
	"net/http"
	"time"

	"github.com/NessibeliY/API/internal/dto"
	"github.com/NessibeliY/API/internal/validator"
	"github.com/gin-gonic/gin"
)

func (c *Client) CreateClientDocument(ctx *gin.Context) {
	date := time.Now().UTC()
	userEmail, err := ctx.Cookie("user") // TODO "user" move to value/contants
	baseResponse := dto.BaseResponse{
		Message: "failed to extract cookie \"user\" from request",
		Status:  http.StatusUnauthorized,
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

	err = c.DocumentServices.CreateDocument(&request, date, userEmail)
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
			Message: "document successfully created", // TODO move messages to value/contants
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
		ctx.Error(err)
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
		ctx.Error(err) // TODO ctx.Error middleware that processes gin errors that reads status from Error
		return
	}

	document, err := c.DocumentServices.GetDocument(&request)
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
