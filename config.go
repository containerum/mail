package main

import (
	"errors"

	"git.containerum.net/ch/mail-templater/storages"
	"git.containerum.net/ch/mail-templater/upstreams"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/mailgun/mailgun-go.v1"
)

func setupLogger() error {
	switch gin.Mode() {
	case gin.TestMode, gin.DebugMode:
	case gin.ReleaseMode:
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
	viper.SetDefault("log_level", logrus.InfoLevel)
	level := logrus.Level(viper.GetInt("log_level"))
	if level > logrus.DebugLevel || level < logrus.PanicLevel {
		return errors.New("invalid log level")
	}
	return nil
}

func getTemplatesStorage() (*storages.TemplateStorage, error) {
	viper.SetDefault("template_db", "templates.db")
	file := viper.GetString("template_db")
	return storages.NewTemplateStorage(file, nil)
}

func getMessagesStorage() (*storages.MessagesStorage, error) {
	viper.SetDefault("messages_db", "messages.db")
	file := viper.GetString("messages_db")
	return storages.NewMessagesStorage(file, nil)
}

func getUpstream(msgStorage *storages.MessagesStorage) (upstreams.Upstream, error) {
	viper.SetDefault("upstream", "mailgun")
	switch viper.GetString("upstream") {
	case "mailgun":
		mg, err := mailgun.NewMailgunFromEnv()
		if err != nil {
			return nil, err
		}
		return upstreams.NewMailgun(mg, msgStorage, viper.GetString("sender_name"), viper.GetString("sender_mail")), nil
	default:
		return nil, errors.New("invalid upstream")
	}
}

func getListenAddr() string {
	viper.SetDefault("listen_addr", ":7070")
	return viper.GetString("listen_addr")
}
