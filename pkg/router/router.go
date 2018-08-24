package router

import (
	"time"

	"net/http"

	"git.containerum.net/ch/auth/static"
	"git.containerum.net/ch/mail-templater/pkg/mterrors"
	h "git.containerum.net/ch/mail-templater/pkg/router/handlers"
	m "git.containerum.net/ch/mail-templater/pkg/router/middleware"
	"github.com/containerum/utils/httputil"
	"github.com/gin-gonic/contrib/ginrus"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	cors "gopkg.in/gin-contrib/cors.v1"

	"github.com/containerum/cherry/adaptors/cherrylog"
	"github.com/containerum/cherry/adaptors/gonic"
	"github.com/containerum/kube-client/pkg/model"
)

//CreateRouter initialises router and middlewares
func CreateRouter(svc *m.Services, status *model.ServiceStatus, enableCORS bool) http.Handler {
	e := gin.New()
	initMiddlewares(e, svc)
	initRoutes(e, status, enableCORS)
	return e
}

func initMiddlewares(e *gin.Engine, svc *m.Services) {
	/* System */
	e.Use(ginrus.Ginrus(logrus.WithField("component", "gin"), time.RFC3339, true))
	e.Use(gonic.Recovery(mterrors.ErrInternalError, cherrylog.NewLogrusAdapter(logrus.WithField("component", "gin"))))
	/* Custom */
	e.Use(m.RegisterServices(svc))
}

// Setup sets up routes
func initRoutes(e *gin.Engine, status *model.ServiceStatus, enableCORS bool) {
	if enableCORS {
		cfg := cors.DefaultConfig()
		cfg.AllowAllOrigins = true
		cfg.AddAllowMethods(http.MethodDelete)
		cfg.AddAllowHeaders(httputil.UserRoleXHeader)
		e.Use(cors.New(cfg))
	}
	e.Group("/static").
		StaticFS("/", static.HTTP)

	e.GET("/status", httputil.ServiceStatus(status))

	e.POST("/send", m.CheckActive(), h.SimpleSendHandler)

	messages := e.Group("/messages")
	{
		messages.GET("/:message_id", m.RequireAdminRole, h.MessageGetHandler)
		messages.GET("", m.RequireAdminRole, h.MessageListGetHandler)
	}

	templates := e.Group("/templates")
	{
		templates.GET("", m.RequireAdminRole, h.TemplateListGetHandler)
		templates.POST("", m.RequireAdminRole, h.TemplateCreateHandler)
		templates.GET("/:name", m.RequireAdminRole, h.TemplateGetHandler)
		templates.PUT("/:name", m.RequireAdminRole, h.TemplateUpdateHandler)
		templates.DELETE("/:name", m.RequireAdminRole, h.TemplateDeleteHandler)
		templates.POST("/:name", h.SendHandler)
	}
}
