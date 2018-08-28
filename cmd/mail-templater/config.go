package main

import (
	"errors"

	"io/ioutil"

	"strings"

	"git.containerum.net/ch/mail-templater/pkg/clients"
	"git.containerum.net/ch/mail-templater/pkg/models"
	"git.containerum.net/ch/mail-templater/pkg/mterrors"
	"git.containerum.net/ch/mail-templater/pkg/storages"
	"git.containerum.net/ch/mail-templater/pkg/storages/bolt"
	"git.containerum.net/ch/mail-templater/pkg/upstreams"
	"github.com/containerum/cherry"
	"github.com/gin-gonic/gin"
	"github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"gopkg.in/mailgun/mailgun-go.v1"
)

const (
	portFlag             = "port"
	debugFlag            = "debug"
	textlogFlag          = "textlog"
	templatesStorageFlag = "template_storage"
	templatesDBFlag      = "template_db"
	messagesStorageFlag  = "messages_storage"
	messagesDBFlag       = "messages_db"
	userManagerFlag      = "um"
	userManagerURLFlag   = "um_url"
	upstreamFlag         = "upstream"
	senderNameFlag       = "sender_name"
	senderMailFlag       = "sender_mail"
	upstreamSimpleFlag   = "upstream_simple"
	senderNameSimpleFlag = "sender_name_simple"
	senderMailSimpleFlag = "sender_mail_simple"
	smtpAddrFlag         = "smtp_addr"
	smtpLoginFlag        = "smtp_login"
	smtpPasswordFlag     = "smtp_password"
	corsFlag             = "cors"
	defaultTemplatesFlag = "default_templates"
)

var flags = []cli.Flag{
	cli.StringFlag{
		EnvVar: "CH_MAIL_PORT",
		Name:   portFlag,
		Value:  "7070",
		Usage:  "port for solutions server",
	},
	cli.BoolFlag{
		EnvVar: "CH_MAIL_DEBUG",
		Name:   debugFlag,
		Usage:  "start the server in debug mode",
	},
	cli.BoolFlag{
		EnvVar: "CH_MAIL_TEXTLOG",
		Name:   textlogFlag,
		Usage:  "output log in text format",
	},
	cli.StringFlag{
		EnvVar: "CH_MAIL_TEMPLATE_STORAGE",
		Name:   templatesStorageFlag,
		Value:  "bolt",
		Usage:  "Templates storage",
	},
	cli.StringFlag{
		EnvVar: "CH_MAIL_TEMPLATE_DB",
		Name:   templatesDBFlag,
		Value:  "templates.db",
		Usage:  "Templates db",
	},
	cli.StringFlag{
		EnvVar: "CH_MAIL_MESSAGES_STORAGE",
		Name:   messagesStorageFlag,
		Value:  "bolt",
		Usage:  "Messages storage",
	},
	cli.StringFlag{
		EnvVar: "CH_MAIL_MESSAGES_DB",
		Name:   messagesDBFlag,
		Value:  "messages.db",
		Usage:  "Messages db",
	},
	cli.StringFlag{
		EnvVar: "CH_MAIL_USER_MANAGER",
		Name:   userManagerFlag,
		Value:  "http",
		Usage:  "User manager kind",
	},
	cli.StringFlag{
		EnvVar: "CH_MAIL_USER_MANAGER_URL",
		Name:   userManagerURLFlag,
		Value:  "http://user-manager:8111",
		Usage:  "User manager",
	},
	cli.StringFlag{
		EnvVar: "CH_MAIL_UPSTREAM",
		Name:   upstreamFlag,
		Value:  "mailgun",
		Usage:  "Upstream (SMTP, Mailgun)",
	},
	cli.StringFlag{
		EnvVar: "CH_MAIL_SENDER_NAME",
		Name:   senderNameFlag,
		Usage:  "Sender name",
	},
	cli.StringFlag{
		EnvVar: "CH_MAIL_SENDER_MAIL",
		Name:   senderMailFlag,
		Usage:  "Sender email",
	},
	cli.StringFlag{
		EnvVar: "CH_MAIL_UPSTREAM_SIMPLE",
		Name:   upstreamSimpleFlag,
		Value:  "mailgun",
		Usage:  "Upstream for simple send method (SMTP, Mailgun)",
	},
	cli.StringFlag{
		EnvVar: "CH_MAIL_SENDER_NAME_SIMPLE",
		Name:   senderNameSimpleFlag,
		Usage:  "Sender name for simple send method",
	},
	cli.StringFlag{
		EnvVar: "CH_MAIL_SENDER_MAIL_SIMPLE",
		Name:   senderMailSimpleFlag,
		Usage:  "Sender email for simple send method",
	},
	cli.StringFlag{
		EnvVar: "CH_MAIL_SMTP_ADDR",
		Name:   smtpAddrFlag,
		Usage:  "Sender email for simple send method",
	},
	cli.StringFlag{
		EnvVar: "CH_MAIL_SMTP_LOGIN",
		Name:   smtpLoginFlag,
		Usage:  "SMTP login",
	},
	cli.StringFlag{
		EnvVar: "CH_MAIL_SMTP_PASSWORD",
		Name:   smtpPasswordFlag,
		Usage:  "SMTP password",
	},
	cli.BoolFlag{
		EnvVar: "CH_MAIL_CORS",
		Name:   corsFlag,
		Usage:  "enable CORS",
	},
	cli.StringFlag{
		EnvVar: "CH_MAIL_DEFAULT_TEMPLATES",
		Name:   defaultTemplatesFlag,
		Value:  "templates.json",
		Usage:  "json file with default templates",
	},
}

func setupLogs(c *cli.Context) {
	if c.Bool("debug") {
		gin.SetMode(gin.DebugMode)
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		gin.SetMode(gin.ReleaseMode)
		logrus.SetLevel(logrus.InfoLevel)
	}

	if c.Bool(textlogFlag) {
		logrus.SetFormatter(&logrus.TextFormatter{})
	} else {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
}

func getTemplatesStorage(c *cli.Context) (storages.TemplateStorage, error) {
	switch c.String(templatesStorageFlag) {
	case "bolt":
		return bolt.NewBoltTemplateStorage(c.String(templatesDBFlag), nil)
	default:
		return nil, errors.New("invalid template storage")
	}
}

func getMessagesStorage(c *cli.Context) (storages.MessagesStorage, error) {
	switch c.String(messagesStorageFlag) {
	case "bolt":
		return bolt.NewBoltMessagesStorage(c.String(messagesDBFlag), nil)
	default:
		return nil, errors.New("invalid messages storage")
	}
}

func getUpstream(c *cli.Context, msgStorage storages.MessagesStorage) (upstreams.Upstream, error) {
	switch c.String(upstreamFlag) {
	case "mailgun":
		mg, err := mailgun.NewMailgunFromEnv()
		if err != nil {
			return nil, err
		}
		return upstreams.NewMailgun(mg, msgStorage, c.String(senderNameFlag), c.String(senderMailFlag)), nil
	case "smtp":
		upstream := upstreams.NewSMTPUpstream(msgStorage, c.String(senderNameFlag), c.String(senderMailFlag), c.String(smtpAddrFlag), c.String(smtpLoginFlag), c.String(smtpPasswordFlag))
		return upstream, nil
	case "dummy":
		return upstreams.NewDummyUpstream(), nil
	default:
		return nil, errors.New("invalid upstream")
	}
}

func getUpstreamSimple(c *cli.Context, msgStorage storages.MessagesStorage) (upstreams.Upstream, error) {
	switch c.String(upstreamSimpleFlag) {
	case "mailgun":
		mg, err := mailgun.NewMailgunFromEnv()
		if err != nil {
			return nil, err
		}
		return upstreams.NewMailgun(mg, msgStorage, c.String(senderNameSimpleFlag), c.String(senderMailSimpleFlag)), nil
	case "smtp":
		upstream := upstreams.NewSMTPUpstream(msgStorage, c.String(senderNameSimpleFlag), c.String(senderMailSimpleFlag), c.String(smtpAddrFlag), c.String(smtpLoginFlag), c.String(smtpPasswordFlag))
		err := upstream.CheckStatus()
		return upstream, err
	case "dummy":
		return upstreams.NewDummyUpstream(), nil
	default:
		return nil, errors.New("invalid upstream")
	}
}

func getUserManagerClient(c *cli.Context) (clients.UserManagerClient, error) {
	switch c.String(userManagerFlag) {
	case "http":
		return clients.NewHTTPUserManagerClient(c.String(userManagerURLFlag)), nil
	default:
		return nil, errors.New("invalid user manager client")
	}
}

func importDefaultTemplates(c *cli.Context, tstorage storages.TemplateStorage) error {
	file, err := ioutil.ReadFile(c.String(defaultTemplatesFlag))
	if err != nil {
		return err
	}
	var defTemplates []models.Template
	if err := jsoniter.Unmarshal(file, &defTemplates); err != nil {
		return err
	}
	var tmplErrs []string
	for _, tmpl := range defTemplates {
		if err := tstorage.PutTemplate(tmpl.Name, tmpl.Version, tmpl.Data, tmpl.Subject, true); err != nil {
			//If template with this name already exists it's not a problem
			if !cherry.Equals(err, mterrors.ErrTemplateAlreadyExists()) {
				logrus.
					WithField("name", tmpl.Name).
					WithField("version", tmpl.Version).
					WithError(err).
					Warn("Unable to import template")

				tmplErrs = append(tmplErrs, err.Error())
			}
		}
	}
	logrus.Println("Templates import finished")
	if len(tmplErrs) != 0 {
		return errors.New(strings.Join(tmplErrs, "; "))
	}
	return nil
}
