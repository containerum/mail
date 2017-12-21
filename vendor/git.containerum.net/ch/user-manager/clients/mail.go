package clients

import (
	"git.containerum.net/ch/mail-templater/upstreams"
	"git.containerum.net/ch/utils"
	"github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"
)

type MailClient struct {
	rest *resty.Client
	log  *logrus.Entry
}

func NewMailClient(serverUrl string) *MailClient {
	log := logrus.WithField("component", "mail_client")
	client := resty.New().SetHostURL(serverUrl).SetLogger(log.WriterLevel(logrus.DebugLevel))
	return &MailClient{
		rest: client,
		log:  log,
	}
}

func (mc *MailClient) sendOneTemplate(tmplName string, recipient *upstreams.Recipient) error {
	req := &upstreams.SendRequest{}
	req.Delay = 0
	req.Message.Recipients = append(req.Message.Recipients, *recipient)
	resp, err := mc.rest.R().
		SetBody(req).
		SetResult(upstreams.SendResponse{}).
		SetError(utils.Error{}).
		Post("/templates/" + tmplName)
	if err != nil {
		return err
	}
	return resp.Error().(*utils.Error)
}

func (mc *MailClient) SendConfirmationMail(recipient *upstreams.Recipient) error {
	mc.log.Infoln("Sending confirmation mail to", recipient.Email)
	return mc.sendOneTemplate("confirm_reg", recipient)
}

func (mc *MailClient) SendActivationMail(recipient *upstreams.Recipient) error {
	mc.log.Infoln("Sending confirmation mail to", recipient.Email)
	return mc.sendOneTemplate("activate_acc", recipient)
}

func (mc *MailClient) SendBlockedMail(recipient *upstreams.Recipient) error {
	mc.log.Infoln("Sending blocked mail to", recipient.Email)
	return mc.sendOneTemplate("blocked_acc", recipient)
}

func (mc *MailClient) SendPasswordChangedMail(recipient *upstreams.Recipient) error {
	mc.log.Infoln("Sending password changed mail to", recipient.Email)
	return mc.sendOneTemplate("pwd_changed", recipient)
}

func (mc *MailClient) SendPasswordResetMail(recipient *upstreams.Recipient) error {
	mc.log.Infoln("Sending reset password mail to", recipient.Email)
	return mc.sendOneTemplate("reset_pwd", recipient)
}
