package client

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/NessibeliY/API/internal/models"
	"github.com/NessibeliY/API/pkg"
	"github.com/gin-gonic/gin"
)

func (c *Client) CORSMiddleware() gin.HandlerFunc { // TODO gin CORS
	allowedOrigins := []string{"http://google.com", "http://facebook.com"}
	return func(ctx *gin.Context) {
		origin := ctx.Request.Header.Get("Origin")

		if pkg.Contains(allowedOrigins, origin) {
			ctx.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.Next()
	}
}

func (c *Client) BasicAuthMiddleware() gin.HandlerFunc {
	unauthorized := func(ctx *gin.Context) {
		ctx.Header("WWW-Authenticate", `Basic realm="Restricted"`)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": false, "message": "user not authorized"})
	}
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			unauthorized(ctx)
			return
		}

		if !strings.HasPrefix(authHeader, basicPrefic) {
			unauthorized(ctx)
			return
		}

		payload := strings.TrimPrefix(authHeader, basicPrefic)
		decoded, err := base64.StdEncoding.DecodeString(payload)
		if err != nil {
			unauthorized(ctx)
			return
		}

		pair := strings.Split(string(decoded), ":")
		if len(pair) != 2 {
			unauthorized(ctx)
			return
		}

		username := pair[0]
		password := pair[1]

		// Check if username and password are not empty
		if username == "" || password == "" {
			unauthorized(ctx)
			return
		}

		// Store the username and password in the context
		ctx.Set("username", username)
		ctx.Set("password", password)

		ctx.Next()
	}
}

func (c *Client) SessionAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookie, err := ctx.Cookie("session-id")
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		fmt.Println("SessionID in middleware", cookie)

		sessionUser := models.SessionUserClient{}
		err = c.sessionServices.GetSession(cookie, &sessionUser)
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
