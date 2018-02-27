package router

import (
	"time"

	"net/http"

	umtypes "git.containerum.net/ch/json-types/user-manager"
	h "git.containerum.net/ch/mail-templater/pkg/router/handlers"
	m "git.containerum.net/ch/mail-templater/pkg/router/middleware"
	"github.com/gin-gonic/contrib/ginrus"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func CreateRouter(svc *m.Services) http.Handler {
	e := gin.New()
	initMiddlewares(e, svc)
	initRoutes(e)
	return e
}

func initMiddlewares(e *gin.Engine, svc *m.Services) {
	/* System */
	e.Use(ginrus.Ginrus(logrus.WithField("component", "gin"), time.RFC3339, true))
	e.Use(gin.RecoveryWithWriter(logrus.WithField("component", "gin_recovery").WriterLevel(logrus.ErrorLevel)))
	/* Custom */
	e.Use(m.RegisterServices(svc))
}

// Setup sets up routes
func initRoutes(e *gin.Engine) {

	requireIdentityHeaders := m.RequireHeaders(umtypes.UserIDHeader, umtypes.UserRoleHeader, umtypes.SessionIDHeader)

	e.POST("/send", h.SimpleSendHandler)

	messages := e.Group("/messages")
	{
		messages.GET("/:message_id", requireIdentityHeaders, m.RequireAdminRole, h.MessageGetHandler)
		messages.GET("/", requireIdentityHeaders, m.RequireAdminRole, h.MessageListGetHandler)
	}

	templates := e.Group("/templates")
	{
		templates.GET("/", requireIdentityHeaders, m.RequireAdminRole, h.TemplateListGetHandler)
		templates.POST("/", requireIdentityHeaders, m.RequireAdminRole, h.TemplateCreateHandler)
		templates.GET("/:name", requireIdentityHeaders, m.RequireAdminRole, h.TemplateGetHandler)
		templates.PUT("/:name", requireIdentityHeaders, m.RequireAdminRole, h.TemplateUpdateHandler)
		templates.DELETE("/:name", requireIdentityHeaders, m.RequireAdminRole, h.TemplateDeleteHandler)
		templates.POST("/:name", h.TemplateSendHandler)
	}
}
