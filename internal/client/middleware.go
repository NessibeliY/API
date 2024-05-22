package client

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/NessibeliY/API/internal/models"
	"github.com/gin-gonic/gin"
)

// TODO move all middlewares here
func (c *Client) BasicAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.Header("WWW-Authenticate", `Basic realm="Restricted"`)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": false, "message": "user not authorized"})
			return
		}

		if !strings.HasPrefix(authHeader, basicPrefic) {
			ctx.Header("WWW-Authenticate", `Basic realm="Restricted"`) // TODO вынести повторения ошибок
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
