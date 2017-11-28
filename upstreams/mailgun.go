package upstreams

import (
	"bytes"
	"errors"
	"html/template"
	"strings"
	"sync"
	"time"

	"encoding/base64"

	"bitbucket.org/exonch/ch-mail-templater/storages"
	"github.com/sirupsen/logrus"
	"gopkg.in/mailgun/mailgun-go.v1"
)

type Mailgun struct {
	api        mailgun.Mailgun
	log        *logrus.Logger
	msgStorage *storages.MessagesStorage
	senderName string
	senderMail string
}

func NewMailgun(conn mailgun.Mailgun, msgStorage *storages.MessagesStorage, senderName, senderMail string) *Mailgun {
	return &Mailgun{
		api:        conn,
		log:        logrus.WithField("component", "mailgun").Logger,
		msgStorage: msgStorage,
		senderName: senderName,
		senderMail: senderMail,
	}
}

func (mg *Mailgun) executeTemplate(tmpl *template.Template, recipient Recipient, commonVars map[string]string) (string, error) {
	var buf bytes.Buffer
	tmplData := make(map[string]interface{})
	for k, v := range commonVars {
		tmplData[k] = v
	}
	for k, v := range recipient.Variables {
		tmplData[k] = v
	}
	e := mg.log.WithField("recipient", recipient.Email).WithFields(tmplData)
	e.Debug("Executing template")
	if err := tmpl.Execute(&buf, tmplData); err != nil {
		e.WithError(err).Error("Execute template failed")
		return "", err
	}
	return buf.String(), nil
}

func (mg *Mailgun) constructMessage(text, subj, to string, delayMinutes int) *mailgun.Message {
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

func (mg *Mailgun) errCollector(ch chan error, errs *[]string) {
	go func() {
		mu := &sync.Mutex{}
		for err := range ch {
			if err != nil {
				mu.Lock()
				*errs = append(*errs, err.Error())
				mu.Unlock()
			}
		}
	}()
}

func (mg *Mailgun) statusCollector(ch chan SendStatus, statuses *[]SendStatus) {
	go func() {
		mu := &sync.Mutex{}
		for s := range ch {
			mu.Lock()
			*statuses = append(*statuses, s)
			mu.Unlock()
		}
	}()
}

func (mg *Mailgun) Send(templateName string, tsv *storages.TemplateStorageValue, request *SendRequest) (resp *SendResponse, err error) {
	mg.log.Debug("Parsing template")
	templateText, err := base64.StdEncoding.DecodeString(tsv.Data)
	if err != nil {
		mg.log.WithError(err).Error("Template data decode failed")
		return nil, err
	}
	tmpl, err := template.New(templateName).Parse(string(templateText))
	if err != nil {
		mg.log.WithError(err).Debug("Template parsing failed")
		return nil, err
	}

	resp = &SendResponse{}

	var errs []string
	errChan := make(chan error)
	mg.errCollector(errChan, &errs)

	statusChan := make(chan SendStatus)
	mg.statusCollector(statusChan, &resp.Statuses)

	wg := &sync.WaitGroup{}
	wg.Add(len(request.Message.Recipients))
	for _, recipient := range request.Message.Recipients {
		text, err := mg.executeTemplate(tmpl, recipient, request.Message.CommonVariables)
		if err != nil {
			errChan <- err
		}

		msg := mg.constructMessage(text, tsv.Subject, recipient.Email, request.Delay)

		go func(msg *mailgun.Message, recipient Recipient, text string) {
			status, id, err := mg.api.Send(msg)
			if err != nil {
				logrus.WithError(err).Error("Message send failed")
				errChan <- err
				return
			}
			mg.log.WithField("status", status).WithField("id", id).Info("Message sent")
			statusChan <- SendStatus{
				RecipientID:  recipient.ID,
				TemplateName: templateName,
				Status:       status,
			}
			mg.msgStorage.PutValue(id, &storages.MessagesStorageValue{
				UserId:       recipient.ID,
				TemplateName: templateName,
				Variables:    recipient.Variables,
				CreatedAt:    time.Now().UTC(),
				Message:      base64.StdEncoding.EncodeToString([]byte(text)),
			})
			wg.Done()
		}(msg, recipient, text)
	}

	wg.Wait()
	close(errChan)
	close(statusChan)

	if len(errs) > 0 {
		err = errors.New(strings.Join(errs, "; "))
	}
	return resp, err
}
