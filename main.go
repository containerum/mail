package main

import (
	"fmt"
	"os"
	"time"

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

	app := gin.New()
	app.Use(gin.Recovery())
	app.Use(ginrus.Ginrus(logrus.StandardLogger(), time.RFC3339, true))

	ts, err := getTemplatesStorage()
	exitOnErr(err)
	ms, err := getMessagesStorage()
	exitOnErr(err)
	us, err := getUpstream(ms)
	exitOnErr(err)

	routes.Setup(app, &routes.Services{
		TemplateStorage: ts,
		MessagesStorage: ms,
		Upstream:        us,
	})
	exitOnErr(app.Run(getListenAddr()))
}
