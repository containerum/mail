package upstreams

import (
	"context"

	"git.containerum.net/ch/mail-templater/pkg/models"
	"github.com/sirupsen/logrus"
)

type dummyUpstream struct {
	log *logrus.Entry
}

// NewDummyUpstream returns a new dummy email upstream.
// It actually does nothing, only logs actions
func NewDummyUpstream() Upstream {
	return &dummyUpstream{
		log: logrus.WithField("component", "dummy_upstream"),
	}
}

//Send
//Sends dummy email
func (du *dummyUpstream) Send(ctx context.Context, templateName string, tsv *models.Template, request *models.SendRequest) (resp *models.SendResponse, err error) {
	resp = &models.SendResponse{}
	for _, recipient := range request.Message.Recipients {
		du.log.WithField("template", templateName).WithFields(recipient.Variables).Infoln("Sending email to", recipient.Email)
		resp.Statuses = append(resp.Statuses, models.SendStatus{
			RecipientID:  recipient.ID,
			TemplateName: templateName,
			Status:       "OK",
		})
	}
	return
}

//SimpleSend
//Sends dummy email in simple way
func (du *dummyUpstream) SimpleSend(ctx context.Context, templateName string, tsv *models.Template, recipient *models.Recipient) (status *models.SendStatus, err error) {
	du.log.WithField("template", templateName).WithFields(recipient.Variables).Infoln("Sending email to", recipient.Email)
	status = &models.SendStatus{
		RecipientID:  recipient.ID,
		TemplateName: templateName,
		Status:       "OK",
	}
	return
}
