package main

import (
	"errors"

	"git.containerum.net/ch/mail-templater/pkg/clients"
	"git.containerum.net/ch/mail-templater/pkg/storages"
	"git.containerum.net/ch/mail-templater/pkg/storages/bolt"
	"git.containerum.net/ch/mail-templater/pkg/upstreams"
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
	logrus.SetLevel(level)
	return nil
}

func getTemplatesStorage() (storages.TemplateStorage, error) {
	viper.SetDefault("template_storage", "bolt")
	switch viper.GetString("template_storage") {
	case "bolt":
		viper.SetDefault("template_db", "templates.db")
		return bolt.NewBoltTemplateStorage(viper.GetString("template_db"), nil)
	default:
		return nil, errors.New("invalid template storage")
	}
}

func getMessagesStorage() (storages.MessagesStorage, error) {
	viper.SetDefault("messages_storage", "bolt")
	switch viper.GetString("messages_storage") {
	case "bolt":
		viper.SetDefault("messages_db", "messages.db")
		return bolt.NewBoltMessagesStorage(viper.GetString("messages_db"), nil)
	default:
		return nil, errors.New("invalid messages storage")
	}
}

func getUpstream(msgStorage storages.MessagesStorage) (upstreams.Upstream, error) {
	viper.SetDefault("upstream", "mailgun")
	switch viper.GetString("upstream") {
	case "mailgun":
		mg, err := mailgun.NewMailgunFromEnv()
		if err != nil {
			return nil, err
		}
		return upstreams.NewMailgun(mg, msgStorage, viper.GetString("sender_name"), viper.GetString("sender_mail")), nil
	case "smtp":
		return upstreams.NewSMTPUpstream(msgStorage, viper.GetString("sender_name"), viper.GetString("sender_mail"), viper.GetString("smtp_addr"), viper.GetString("smtp_login"), viper.GetString("smtp_password")), nil
	case "dummy":
		return upstreams.NewDummyUpstream(), nil
	default:
		return nil, errors.New("invalid upstream")
	}
}

func getUpstreamSimple(msgStorage storages.MessagesStorage) (upstreams.Upstream, error) {
	viper.SetDefault("upstream_simple", "mailgun")
	switch viper.GetString("upstream_simple") {
	case "mailgun":
		mg, err := mailgun.NewMailgunFromEnv()
		if err != nil {
			return nil, err
		}
		return upstreams.NewMailgun(mg, msgStorage, viper.GetString("sender_name"), viper.GetString("sender_mail")), nil
	case "smtp":
		return upstreams.NewSMTPUpstream(msgStorage, viper.GetString("sender_name"), viper.GetString("sender_mail"), viper.GetString("smtp_addr"), viper.GetString("smtp_login"), viper.GetString("smtp_password")), nil
	case "dummy":
		return upstreams.NewDummyUpstream(), nil
	default:
		return nil, errors.New("invalid upstream")
	}
}

func getListenAddr() string {
	viper.SetDefault("listen_addr", ":7070")
	return viper.GetString("listen_addr")
}

func getUserManagerClient() (clients.UserManagerClient, error) {
	viper.SetDefault("user_manager", "http")
	switch viper.GetString("user_manager") {
	case "http":
		viper.SetDefault("user_manager_url", "http://user-manager:8111")
		return clients.NewHTTPUserManagerClient(viper.GetString("user_manager_url")), nil
	default:
		return nil, errors.New("invalid user manager client")
	}
}