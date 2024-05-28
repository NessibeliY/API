package client

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func (c *Client) Routes(router *gin.Engine) {
	// Apply the CORS middleware
	router.Use(c.CORSMiddleware())

	// Apply the sessions middleware to the router
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("session123", store))

	// Public routes
	router.POST("signup", c.Signup)

	// Login route with Basic Auth Middleware
	router.POST("login", c.BasicAuthMiddleware(), c.Login)

	router.POST("show-document", c.ShowClientDocument)

	// Protected routes
	protected := router.Group("/protected")
	protected.Use(c.SessionAuthMiddleware())
	{
		protected.POST("/create-document", c.CreateClientDocument)
		protected.GET("/dashboard", c.Dashboard)
	}
}
