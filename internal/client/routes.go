package client

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func (c *Client) Routes(router *gin.Engine) {
	// Apply the sessions middleware to the router
	// router.Use(sessions.Sessions("my_session", rdb)) // TODO rename my_session to session?
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

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
