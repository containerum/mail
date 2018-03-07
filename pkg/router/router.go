package router

import (
	"time"

	"net/http"

	h "git.containerum.net/ch/mail-templater/pkg/router/handlers"
	m "git.containerum.net/ch/mail-templater/pkg/router/middleware"
	"github.com/gin-gonic/contrib/ginrus"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	ch "git.containerum.net/ch/kube-client/pkg/cherry"
	"git.containerum.net/ch/kube-client/pkg/cherry/adaptors/cherrylog"
	"git.containerum.net/ch/kube-client/pkg/cherry/adaptors/gonic"
	cherry "git.containerum.net/ch/kube-client/pkg/cherry/mail-templater"
)

//CreateRouter initialises router and middlewares
func CreateRouter(svc *m.Services) http.Handler {
	e := gin.New()
	initMiddlewares(e, svc)
	initRoutes(e)
	return e
}

func initMiddlewares(e *gin.Engine, svc *m.Services) {
	/* System */
	e.Use(ginrus.Ginrus(logrus.WithField("component", "gin"), time.RFC3339, true))
	e.Use(gonic.Recovery(func() *ch.Err { return cherry.ErrInternalError() }, cherrylog.NewLogrusAdapter(logrus.WithField("component", "gin"))))
	/* Custom */
	e.Use(m.RegisterServices(svc))
}

// Setup sets up routes
func initRoutes(e *gin.Engine) {

	e.POST("/send", h.SimpleSendHandler)

	messages := e.Group("/messages")
	{
		messages.GET("/:message_id", m.RequireAdminRole, h.MessageGetHandler)
		messages.GET("/", m.RequireAdminRole, h.MessageListGetHandler)
	}

	templates := e.Group("/templates")
	{
		templates.GET("/", m.RequireAdminRole, h.TemplateListGetHandler)
		templates.POST("/", m.RequireAdminRole, h.TemplateCreateHandler)
		templates.GET("/:name", m.RequireAdminRole, h.TemplateGetHandler)
		templates.PUT("/:name", m.RequireAdminRole, h.TemplateUpdateHandler)
		templates.DELETE("/:name", m.RequireAdminRole, h.TemplateDeleteHandler)
		templates.POST("/:name", h.SendHandler)
	}
}
