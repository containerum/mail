package upstreams

import (
	"bytes"
	"errors"
	"html/template"
	"strings"
	"sync"
	"time"

	"encoding/base64"

	"context"

	mttypes "git.containerum.net/ch/json-types/mail-templater"
	"git.containerum.net/ch/mail-templater/pkg/storages"
	"github.com/sirupsen/logrus"
	"gopkg.in/mailgun/mailgun-go.v1"
)

type mgUpstream struct {
	api        mailgun.Mailgun
	log        *logrus.Entry
	msgStorage storages.MessagesStorage
	senderName string
	senderMail string
}

// NewMailgun returns mail upstream which uses "Mailgun" service to send emails.
func NewMailgun(conn mailgun.Mailgun, msgStorage storages.MessagesStorage, senderName, senderMail string) Upstream {
	return &mgUpstream{
		api:        conn,
		log:        logrus.WithField("component", "mailgun"),
		msgStorage: msgStorage,
		senderName: senderName,
		senderMail: senderMail,
	}
}

func (mg *mgUpstream) executeTemplate(tmpl *template.Template, recipient *mttypes.Recipient, commonVars map[string]string) (string, error) {
	var buf bytes.Buffer
	tmplData := make(map[string]interface{})
	for k, v := range commonVars {
		tmplData[k] = v
	}
	for k, v := range recipient.Variables {
		tmplData[k] = v
	}
	e := mg.log.WithField("recipient", recipient.Email).WithFields(tmplData)
	e.Debugln("Executing template")
	if err := tmpl.Execute(&buf, tmplData); err != nil {
		e.WithError(err).Errorln("Execute template failed")
		return "", err
	}
	return buf.String(), nil
}

func (mg *mgUpstream) constructMessage(text, subj, to string, delayMinutes int) *mailgun.Message {
	msg := mg.api.NewMessage(mg.senderMail, subj, "", to)
	msg.SetHtml(text)
	if delayMinutes > 0 {
		msg.SetDeliveryTime(time.Now().Add(time.Minute * time.Duration(delayMinutes)))
	}
	msg.SetDKIM(true)
	msg.SetTracking(true)
	msg.SetTrackingClicks(true)
	msg.SetTrackingOpens(true)
	return msg
}

func (mg *mgUpstream) errCollector(ch chan error, errs *[]string, wg *sync.WaitGroup) {
	for err := range ch {
		if err != nil {
			mg.log.WithError(err).Debug("caught error")
			*errs = append(*errs, err.Error())
		}
	}
	wg.Done()
}

func (mg *mgUpstream) statusCollector(ch chan mttypes.SendStatus, statuses *[]mttypes.SendStatus, wg *sync.WaitGroup) {
	for s := range ch {
		mg.log.Debugf("caught status: %#v", s)
		*statuses = append(*statuses, s)
	}
	wg.Done()
}

func (mg *mgUpstream) parseTemplate(templateName string, tsv *mttypes.Template) (tmpl *template.Template, err error) {
	mg.log.Debugln("Parsing template ", templateName)
	templateText, err := base64.StdEncoding.DecodeString(tsv.Data)
	if err != nil {
		mg.log.WithError(err).Errorln("Template data decode failed")
		return nil, err
	}
	tmpl, err = template.New(templateName).Parse(string(templateText))
	if err != nil {
		mg.log.WithError(err).Errorln("Template parse failed")
	}
	return tmpl, err
}

//Send
//Sends email using mailgun
func (mg *mgUpstream) Send(ctx context.Context, templateName string, tsv *mttypes.Template, request *mttypes.SendRequest) (resp *mttypes.SendResponse, err error) {

	tmpl, err := mg.parseTemplate(templateName, tsv)
	if err != nil {
		return nil, err
	}

	resp = &mttypes.SendResponse{}

	mgDoneCh := make(chan struct{}) // for cancelling with context support. Mailgun api has no methods with context

	go func() {
		defer close(mgDoneCh)
		wg := sync.WaitGroup{}
		wg.Add(2) // error and status collectors

		var errs []string
		errChan := make(chan error)
		go mg.errCollector(errChan, &errs, &wg)

		statusChan := make(chan mttypes.SendStatus)
		go mg.statusCollector(statusChan, &resp.Statuses, &wg)

		msgWG := sync.WaitGroup{}
		msgWG.Add(len(request.Message.Recipients))
		for _, recipient := range request.Message.Recipients {
			text, tmplError := mg.executeTemplate(tmpl, &recipient, request.Message.CommonVariables)
			if tmplError != nil {
				errChan <- tmplError
				continue
			}

			msg := mg.constructMessage(text, tsv.Subject, recipient.Email, request.Delay)

			go func(msg *mailgun.Message, recipient mttypes.Recipient, text string) {
				defer msgWG.Done()
				status, id, mgErr := mg.api.Send(msg)
				if mgErr != nil {
					mg.log.WithError(mgErr).Errorln("Message send failed")
					errChan <- mgErr
					return
				}
				mg.log.WithField("status", status).WithField("id", id).Infoln("Message sent")
				statusChan <- mttypes.SendStatus{
					RecipientID:  recipient.ID,
					TemplateName: templateName,
					Status:       status,
				}
				errChan <- mg.msgStorage.PutMessage(id, &mttypes.MessagesStorageValue{
					UserId:       recipient.ID,
					TemplateName: templateName,
					Variables:    recipient.Variables,
					CreatedAt:    time.Now().UTC(),
					Message:      base64.StdEncoding.EncodeToString([]byte(text)),
				})
			}(msg, recipient, text)
		}

		msgWG.Wait()
		close(errChan)
		close(statusChan)
		wg.Wait()
		if len(errs) > 0 {
			err = errors.New(strings.Join(errs, "; "))
		}
	}()

	select {
	case <-ctx.Done():
		mg.log.Info("Operation cancelled")
		err = ctx.Err()
	case <-mgDoneCh:
	}

	return resp, err
}

//SimpleSend
//Sends email using mailgun in simple way
func (mg *mgUpstream) SimpleSend(ctx context.Context, templateName string, tsv *mttypes.Template, recipient *mttypes.Recipient) (status *mttypes.SendStatus, err error) {
	tmpl, err := mg.parseTemplate(templateName, tsv)
	if err != nil {
		return nil, err
	}

	text, err := mg.executeTemplate(tmpl, recipient, nil)
	if err != nil {
		return nil, err
	}

	msg := mg.constructMessage(text, tsv.Subject, recipient.Email, 0)

	mgDoneCh := make(chan struct{})
	go func() {
		defer close(mgDoneCh)

		var s, id string

		s, id, err = mg.api.Send(msg)
		if err != nil {
			mg.log.WithError(err).Errorln("Message send failed")
			return
		}
		mg.log.WithField("status", s).WithField("id", id).Infoln("Message sent")

		status = &mttypes.SendStatus{
			RecipientID:  recipient.ID,
			TemplateName: templateName,
			Status:       s,
		}

		err = mg.msgStorage.PutMessage(id, &mttypes.MessagesStorageValue{
			UserId:       recipient.ID,
			TemplateName: templateName,
			Variables:    recipient.Variables,
			CreatedAt:    time.Now().UTC(),
			Message:      base64.StdEncoding.EncodeToString([]byte(text)),
		})
	}()

	select {
	case <-ctx.Done():
		mg.log.Info("Operation cancelled")
		err = ctx.Err()
	case <-mgDoneCh:
	}

	return status, err
}
