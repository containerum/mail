package routes

import (
	"git.containerum.net/ch/mail-templater/clients"
	"git.containerum.net/ch/mail-templater/storages"
	"git.containerum.net/ch/mail-templater/upstreams"
	"github.com/gin-gonic/gin"
)

// Services is a collection of dependencies to perform server operations
type Services struct {
	MessagesStorage   storages.MessagesStorage
	TemplateStorage   storages.TemplateStorage
	Upstream          upstreams.Upstream
	UpstreamSimple    upstreams.Upstream
	UserManagerClient clients.UserManagerClient
}

var svc *Services

// Setup sets up routes
func Setup(app *gin.Engine, services *Services) {
	svc = services

	app.POST("/send", simpleSendHandler)

	messages := app.Group("/messages")
	{
		messages.GET("/:message_id", messageGetHandler)
	}

	templates := app.Group("/templates")
	{
		templates.GET("/", templateListGetHandler)
		templates.POST("/", templateCreateHandler)
		templates.GET("/:name", templateGetHandler)
		templates.PUT("/:name", templateUpdateHandler)
		templates.DELETE("/:name", templateDeleteHandler)
		templates.POST("/:name", templateSendHandler)
	}
}
