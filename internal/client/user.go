package client

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/NessibeliY/API/internal/dto"
	"github.com/NessibeliY/API/internal/models"
	"github.com/NessibeliY/API/internal/validator"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

const basicPrefic = "Basic " // difference between stack and heap

func (c *Client) Login(ctx *gin.Context) {
	var request dto.LoginRequest
	err := ctx.ShouldBindJSON(&request) // TODO difference between bindjson and shouldbindjson
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.userServices.LoginUser(&request)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		ctx.JSON(http.StatusNotFound, gin.H{"error": "email not found"})
		return
	case err != nil:
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// If authentication is successful, set session
	sessionID := uuid.New().String()
	ctx.SetCookie("session-id", sessionID, 200, "/", "localhost", false, true) //TODO read about "/", "localhost", false, true
	sessionUser := models.SessionUserClient{
		Email:         request.Email,
		Authenticated: true,
	}
	ctx.SetCookie("user", sessionUser.Email, 200, "/", "localhost", false, true) //TODO remove session-id from cookie (use userID)

	err = c.sessionServices.SetSession(sessionID, sessionUser, 100*time.Second)
	co, _ := ctx.Cookie("session-id")
	fmt.Println("SessionID set", co)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot save session to redis"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "user logged in", "session-id": sessionID})
}

// TODO pkg CheckPassword unit test, вопросы к собесам
func (c *Client) Dashboard(ctx *gin.Context) {
	user, err := ctx.Cookie("user")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Welcome to your dashboard",
		"user":    user,
	})
}

func (c *Client) Signup(ctx *gin.Context) {
	var request dto.SignupRequest
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

	err = c.userServices.SignupUser(&request)
	switch {
	case errors.Is(err, models.ErrDuplicateEmail):
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "this email already exists"})
		return
	case err != nil:
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user successfully registered"})
}
