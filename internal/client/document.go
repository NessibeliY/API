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
	var request dto.CreateDocumentRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.documentServices.AddInfoAndCreateDocument(&request, date)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "document successfully created"})
}

func (c *Client) ShowClientDocument(ctx *gin.Context) {
	var request dto.ShowDocumentRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	document, err := c.documentServices.ShowDocument(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"1": gin.H{"message": "document successfully read"},
		"2": gin.H{
			"ID":    document.ID,
			"title": document.Title, "content": document.Content, "image path": document.ImagePath,
			"author ID": document.AuthorID, "date created": document.DateCreated,
		},
	})
}
