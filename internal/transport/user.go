package transport

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/NessibeliY/API/config"
	"github.com/NessibeliY/API/internal/dto"
	"github.com/NessibeliY/API/internal/models"
	redisInternal "github.com/NessibeliY/API/internal/redis"
	"github.com/NessibeliY/API/internal/validator"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type userServices interface {
	SignupUser(*dto.SignupRequest) error
	LoginUser(*dto.LoginRequest) error
}

type Transport struct {
	services userServices
}

func NewTransport(services userServices) *Transport {
	return &Transport{services: services}
}

func (t *Transport) Routes(router *gin.Engine, cfg *config.Config) {
	// Set up session store
	rdb, err := redisInternal.NewRedisClient()
	if err != nil {
		log.Println("Failed to connect to Redis:", err)
		return
	}

	// Apply the sessions middleware to the router
	router.Use(sessions.Sessions("my_session", rdb))

	// Public routes
	router.POST("signup", t.Signup)

	// Login route with Basic Auth Middleware
	router.POST("login", t.BasicAuthMiddleware(), t.Login)

	// Protected routes
	protected := router.Group("/protected")
	protected.Use(t.SessionAuthMiddleware())
	{
		protected.GET("/dashboard", t.Dashboard)
	}
}

func (t *Transport) BasicAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.Header("WWW-Authenticate", `Basic realm="Restricted"`)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": false, "message": "user not authorized"})
			return
		}

		const basicPrefic = "Basic "
		if !strings.HasPrefix(authHeader, basicPrefic) {
			ctx.Header("WWW-Authenticate", `Basic realm="Restricted"`)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": false, "message": "user not authorized"})
			return
		}

		payload := strings.TrimPrefix(authHeader, basicPrefic)
		decoded, err := base64.StdEncoding.DecodeString(payload)
		if err != nil {
			ctx.Header("WWW-Authenticate", `Basic realm="Restricted"`)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": false, "message": "user not authorized"})
			return
		}

		pair := strings.Split(string(decoded), ":")
		if len(pair) != 2 {
			ctx.Header("WWW-Authenticate", `Basic realm="Restricted"`)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": false, "message": "user not authorized"})
			return
		}

		username := pair[0]
		password := pair[1]

		// Check if username and password are not empty
		if username == "" || password == "" {
			ctx.Header("WWW-Authenticate", `Basic realm="Restricted"`)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": false, "message": "user not authorized"})
			return
		}

		// Store the username and password in the context
		ctx.Set("username", username)
		ctx.Set("password", password)

		ctx.Next()
	}
}

func (t *Transport) Login(ctx *gin.Context) {
	var request dto.LoginRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = t.services.LoginUser(&request)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		ctx.JSON(http.StatusNotFound, gin.H{"error": "email not found"})
		return
	case err != nil:
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "user logged in"})

	// If authentication is successful, set session
	sessionID := uuid.New().String()
	ctx.Addcookie("session-id", sessionID, 60)

	sessionUser := models.SessionUserClient{
		Email:         request.Email,
		Authenticated: true,
	}
	err := redisClient.Set(sessionID, sessionUser, 60*time.Second)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot save session to redis"})
		return
	}
}

func (t *Transport) SessionAuthMiddleware(redisClient *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookie, err := ctx.Cookie("session-id")
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		sessionUser := models.SessionUserClient{}
		err = redisClient.Get(cookie.Value, &sessionUser)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if auth := sessionUser.Authenticated; !auth {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.Next()
	}
}

func GetSessionData(client *redis.Client, sessionID string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data, err := client.Get(ctx, sessionKeyPrefix+sessionID).Result()
	if err != nil {
		return "", err
	}
	return data, nil
}

func (t *Transport) Dashboard(ctx *gin.Context) {
	session := sessions.Default(ctx)
	user := session.Get("user")

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Welcome to your dashboard",
		"user":    user,
	})
}

func (t *Transport) Signup(ctx *gin.Context) {
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

	err = t.services.SignupUser(&request)
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
