package router

import (
	umtypes "git.containerum.net/ch/json-types/user-manager"
	"git.containerum.net/ch/mail-templater/pkg/clients"
	h "git.containerum.net/ch/mail-templater/pkg/router/handlers"
	m "git.containerum.net/ch/mail-templater/pkg/router/middleware"
	"git.containerum.net/ch/mail-templater/pkg/storages"
	"git.containerum.net/ch/mail-templater/pkg/upstreams"
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

var Svc *Services

// Setup sets up routes
func Setup(app *gin.Engine, services *Services) {
	Svc = services

	requireIdentityHeaders := m.RequireHeaders(umtypes.UserIDHeader, umtypes.UserRoleHeader, umtypes.SessionIDHeader)

	app.POST("/send", h.SimpleSendHandler)

	messages := app.Group("/messages")
	{
		messages.GET("/:message_id", requireIdentityHeaders, m.RequireAdminRole, h.MessageGetHandler)
		messages.GET("/", requireIdentityHeaders, m.RequireAdminRole, h.MessageListGetHandler)
	}

	templates := app.Group("/templates")
	{
		templates.GET("/", requireIdentityHeaders, m.RequireAdminRole, h.TemplateListGetHandler)
		templates.POST("/", requireIdentityHeaders, m.RequireAdminRole, h.TemplateCreateHandler)
		templates.GET("/:name", requireIdentityHeaders, m.RequireAdminRole, h.TemplateGetHandler)
		templates.PUT("/:name", requireIdentityHeaders, m.RequireAdminRole, h.TemplateUpdateHandler)
		templates.DELETE("/:name", requireIdentityHeaders, m.RequireAdminRole, h.TemplateDeleteHandler)
		templates.POST("/:name", h.TemplateSendHandler)
	}
}
