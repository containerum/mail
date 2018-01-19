package main

import (
	"fmt"
	"os"
	"time"

	"context"
	"net/http"
	"os/signal"

	"git.containerum.net/ch/mail-templater/routes"
	"github.com/gin-gonic/contrib/ginrus"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func exitOnErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	viper.SetEnvPrefix("ch_mail")
	viper.AutomaticEnv()
	exitOnErr(setupLogger())

	logrus.Infoln("starting server...")

	app := gin.New()
	app.Use(gin.RecoveryWithWriter(logrus.StandardLogger().WithField("component", "gin_recovery").WriterLevel(logrus.ErrorLevel)))
	app.Use(ginrus.Ginrus(logrus.StandardLogger(), time.RFC3339, true))

	ts, err := getTemplatesStorage()
	defer ts.Close()
	exitOnErr(err)
	ms, err := getMessagesStorage()
	defer ms.Close()
	exitOnErr(err)
	us, err := getUpstream(ms)
	exitOnErr(err)
	uss, err := getUpstreamSimple(ms)
	exitOnErr(err)
	um, err := getUserManagerClient()
	exitOnErr(err)

	routes.Setup(app, &routes.Services{
		TemplateStorage:   ts,
		MessagesStorage:   ms,
		Upstream:          us,
		UpstreamSimple:    uss,
		UserManagerClient: um,
	})

	// graceful shutdown support

	srv := http.Server{
		Addr:    getListenAddr(),
		Handler: app,
	}

	go exitOnErr(srv.ListenAndServe())

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logrus.Infoln("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	exitOnErr(srv.Shutdown(ctx))
}
