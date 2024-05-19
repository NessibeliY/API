package transport

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

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

const basicPrefic = "Basic " // difference between stack and heap

type userServices interface { // TODO SOLID, пустой интерфейс. как интерфейсы го отличаются от других, утиная типизация, контракт (отличие го от С++), под капотом интерфейса есть тип и дата
	SignupUser(*dto.SignupRequest) error
	LoginUser(*dto.LoginRequest) error
	SetSession(string, models.SessionUserClient, time.Duration) error
	GetSession(string, *models.SessionUserClient) error
}

type Transport struct {
	services userServices
	// cfg *config.Config глобальный
}

func NewTransport(services userServices) *Transport {
	return &Transport{services: services}
}

func (t *Transport) Routes(router *gin.Engine) {
	// Apply the sessions middleware to the router
	// router.Use(sessions.Sessions("my_session", rdb)) // TODO rename my_session to session?
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

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

func (t *Transport) Login(ctx *gin.Context) {
	var request dto.LoginRequest
	err := ctx.ShouldBindJSON(&request) // TODO difference between bindjson and shouldbindjson
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

	// If authentication is successful, set session
	// session := sessions.Default(ctx)
	sessionID := uuid.New().String()
	// session.Set("session-id", sessionID)
	// session.Set("user", request.Email)
	// session.Options(sessions.Options{
	// 	Path:     "/",
	// 	MaxAge:   60 * 1,
	// 	HttpOnly: true,
	// })
	// if err = session.Save(); err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save session"})
	// 	return
	// }
	fmt.Println("SessionID uuid", sessionID)
	ctx.SetCookie("session-id", sessionID, 100, "/", "localhost", false, true)
	sessionUser := models.SessionUserClient{
		Email:         request.Email,
		Authenticated: true,
	}
	ctx.SetCookie("user", sessionUser.Email, 100, "/", "localhost", false, true)

	err = t.services.SetSession(sessionID, sessionUser, 60*time.Second)
	co, _ := ctx.Cookie("session-id")
	fmt.Println("SessionID set", co)
	// fmt.Println("Session ID set:", session.Get("session-id"))
	// fmt.Println("User set:", session.Get("user"))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot save session to redis"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "user logged in", "session-id": sessionID})
}

func (t *Transport) SessionAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// session := sessions.Default(ctx)
		// sessionID := session.Get("session-id")
		// fmt.Println("Session ID in middleware:", sessionID)
		// if sessionID == nil {
		// 	ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "unauthorized"})
		// 	return
		// }
		cookie, err := ctx.Cookie("session-id")
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		fmt.Println("SessionID in middleware", cookie)

		sessionUser := models.SessionUserClient{}
		err = t.services.GetSession(cookie, &sessionUser)
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

// TODO pkg CheckPassword unit test, redis, вопросы к собесам
func (t *Transport) Dashboard(ctx *gin.Context) {
	user, err := ctx.Cookie("user")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

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
