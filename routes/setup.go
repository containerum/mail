package routes

import (
	umtypes "git.containerum.net/ch/json-types/user-manager"
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

	requireIdentityHeaders := requireHeaders(umtypes.UserIDHeader, umtypes.UserRoleHeader, umtypes.SessionIDHeader)

	app.POST("/send", simpleSendHandler)

	messages := app.Group("/messages")
	{
		messages.GET("/:message_id", requireIdentityHeaders, requireAdminRole, messageGetHandler)
	}

	templates := app.Group("/templates")
	{
		templates.GET("/", requireIdentityHeaders, requireAdminRole, templateListGetHandler)
		templates.POST("/", requireIdentityHeaders, requireAdminRole, templateCreateHandler)
		templates.GET("/:name", requireIdentityHeaders, requireAdminRole, templateGetHandler)
		templates.PUT("/:name", requireIdentityHeaders, requireAdminRole, templateUpdateHandler)
		templates.DELETE("/:name", requireIdentityHeaders, requireAdminRole, templateDeleteHandler)
		templates.POST("/:name", templateSendHandler)
	}
}
