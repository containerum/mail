package routes

import (
	"git.containerum.net/ch/mail-templater/clients"
	"git.containerum.net/ch/mail-templater/storages"
	"git.containerum.net/ch/mail-templater/upstreams"
	"github.com/gin-gonic/gin"
)

type Services struct {
	MessagesStorage   *storages.MessagesStorage
	TemplateStorage   *storages.TemplateStorage
	Upstream          upstreams.Upstream
	UserManagerClient clients.UserManagerClient
}

var svc *Services

func Setup(app *gin.Engine, services *Services) {
	svc = services

	app.POST("/send", simpleSendHandler)

	messages := app.Group("/messages")
	{
		messages.GET("/:message_id", messageGetHandler)
	}

	templates := app.Group("/templates")
	{
		templates.POST("/", templateCreateHandler)
		templates.GET("/:name", templateGetHandler)
		templates.PUT("/:name", templateUpdateHandler)
		templates.DELETE("/:name", templateDeleteHandler)
		templates.POST("/:name", templateSendHandler)
	}
}
