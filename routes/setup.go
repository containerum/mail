package routes

import (
	"git.containerum.net/ch/mail-templater/storages"
	"git.containerum.net/ch/mail-templater/upstreams"
	"github.com/gin-gonic/gin"
)

type Services struct {
	MessagesStorage *storages.MessagesStorage
	TemplateStorage *storages.TemplateStorage
	Upstream        upstreams.Upstream
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
		templates.GET("/:template_name", templateGetHandler)
		templates.PUT("/:template_name", templateUpdateHandler)
		templates.DELETE("/:template_name", templateDeleteHandler)
		templates.POST("/:template_name", templateSendHandler)
	}
}
