package upstreams

import (
	"bytes"
	"errors"
	"html/template"
	"strings"
	"sync"
	texttemplate "text/template"
	"time"

	"encoding/base64"

	"context"

	"crypto/tls"
	"math/rand"
	"net/smtp"
	"strconv"

	"net"

	mttypes "git.containerum.net/ch/json-types/mail-templater"
	"git.containerum.net/ch/mail-templater/storages"
	"github.com/sirupsen/logrus"
)

type smtpUpstream struct {
	log          *logrus.Entry
	msgStorage   storages.MessagesStorage
	senderName   string
	senderMail   string
	smtpAddress  string
	smtpLogin    string
	smtpPassword string
}

func NewSmtpUpstream(msgStorage storages.MessagesStorage, senderName string, senderMail string, smtpAddress string, smtpLogin string, smtpPassword string) Upstream {
	return &smtpUpstream{
		log:          logrus.WithField("component", "smtp"),
		msgStorage:   msgStorage,
		senderName:   senderName,
		senderMail:   senderMail,
		smtpAddress:  smtpAddress,
		smtpLogin:    smtpLogin,
		smtpPassword: smtpPassword,
	}
}

type mailData struct {
	SenderName    string
	SenderMail    string
	RecipientName string
	RecipientMail string
	Subject       string
	MessageID     string
	Body          string
}

const emailtemplate = `From: {{.SenderName}} <{{.SenderMail}}>
To: {{.RecipientName}} <{{.RecipientMail}}>
Subject: {{.Subject}}
Message-ID: {{.MessageID}}
MIME-version: 1.0;
Content-Type: text/html; charset="UTF-8";

{{.Body}}`

func (smtpu *smtpUpstream) executeTemplate(tmpl *template.Template, recipient *mttypes.Recipient, commonVars map[string]string) (string, error) {
	var buf bytes.Buffer
	tmplData := make(map[string]interface{})
	for k, v := range commonVars {
		tmplData[k] = v
	}
	for k, v := range recipient.Variables {
		tmplData[k] = v
	}
	e := smtpu.log.WithField("recipient", recipient.Email).WithFields(tmplData)
	e.Debugln("Executing template")
	if err := tmpl.Execute(&buf, tmplData); err != nil {
		e.WithError(err).Errorln("Execute template failed")
		return "", err
	}
	return buf.String(), nil
}

func (smtpu *smtpUpstream) constructMessage(template *texttemplate.Template, recipient mttypes.Recipient, msgID string, subject string, text string) (string, error) {
	newmail := mailData{SenderName: smtpu.senderName,
		SenderMail:    smtpu.senderMail,
		RecipientName: recipient.Name,
		RecipientMail: recipient.Email,
		Subject:       subject,
		MessageID:     msgID,
		Body:          text}

	var mailtext bytes.Buffer
	err := template.Execute(&mailtext, newmail)
	if err != nil {
		return "Error", err
	}

	return mailtext.String(), nil
}

func (smtpu *smtpUpstream) errCollector(ch chan error, errs *[]string, wg *sync.WaitGroup) {
	for err := range ch {
		if err != nil {
			smtpu.log.WithError(err).Debug("caught error")
			*errs = append(*errs, err.Error())
		}
	}
	wg.Done()
}

func (smtpu *smtpUpstream) statusCollector(ch chan mttypes.SendStatus, statuses *[]mttypes.SendStatus, wg *sync.WaitGroup) {
	for s := range ch {
		smtpu.log.Debugf("caught status: %#v", s)
		*statuses = append(*statuses, s)
	}
	wg.Done()
}

func (smtpu *smtpUpstream) parseTemplate(templateName string, tsv *mttypes.TemplateStorageValue) (tmpl *template.Template, err error) {
	smtpu.log.Debugln("Parsing template ", templateName)
	templateText, err := base64.StdEncoding.DecodeString(tsv.Data)
	if err != nil {
		smtpu.log.WithError(err).Errorln("Template data decode failed")
		return nil, err
	}
	tmpl, err = template.New(templateName).Parse(string(templateText))
	if err != nil {
		smtpu.log.WithError(err).Errorln("Template parse failed")
	}
	return tmpl, err
}

func (smtpu *smtpUpstream) Send(ctx context.Context, templateName string, tsv *mttypes.TemplateStorageValue, request *mttypes.SendRequest) (resp *mttypes.SendResponse, err error) {

	tmpl, err := smtpu.parseTemplate(templateName, tsv)
	if err != nil {
		return nil, err
	}

	tmplemail, err := texttemplate.New("emailtemplate").Parse(emailtemplate)
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
		go smtpu.errCollector(errChan, &errs, &wg)

		statusChan := make(chan mttypes.SendStatus)
		go smtpu.statusCollector(statusChan, &resp.Statuses, &wg)

		msgWG := sync.WaitGroup{}
		msgWG.Add(len(request.Message.Recipients))
		for _, recipient := range request.Message.Recipients {
			text, err := smtpu.executeTemplate(tmpl, &recipient, request.Message.CommonVariables)
			if err != nil {
				errChan <- err
				msgWG.Done()
				continue
			}

			msgID := time.Now().UTC().Format("20060102150405.") + strconv.Itoa(rand.Int())
			mailtext, err := smtpu.constructMessage(tmplemail, recipient, msgID, tsv.Subject, text)
			if err != nil {
				errChan <- err
				msgWG.Done()
				continue
			}

			auth := smtp.PlainAuth("", smtpu.smtpLogin, smtpu.smtpPassword, smtpu.smtpAddress)

			tlsconfig := &tls.Config{
				InsecureSkipVerify: true,
				ServerName:         smtpu.smtpAddress,
			}

			conn, err := tls.Dial("tcp", smtpu.smtpAddress, tlsconfig)
			if err != nil {
				smtpu.log.WithError(err).Errorln("Message send failed")
				errChan <- err
				return
			}
			client, err := smtp.NewClient(conn, smtpu.smtpAddress)
			defer client.Quit()
			if err != nil {
				smtpu.log.WithError(err).Errorln("Message send failed")
				errChan <- err
				return
			}
			if err = client.Auth(auth); err != nil {
				smtpu.log.WithError(err).Errorln("Message send failed")
				errChan <- err
				return
			}
			if err = client.Mail(smtpu.senderMail); err != nil {
				smtpu.log.WithError(err).Errorln("Message send failed")
				errChan <- err
				return
			}
			if err = client.Rcpt(recipient.Email); err != nil {
				smtpu.log.WithError(err).Errorln("Message send failed (wrong email)")
				errChan <- err
				msgWG.Done()
				continue
			}

			go func(client *smtp.Client, recipient mttypes.Recipient, text string, messageID string) {
				defer msgWG.Done()
				defer client.Quit()

				w, err := client.Data()
				if err != nil {
					smtpu.log.WithError(err).Errorln("Message send failed")
					errChan <- err
					return
				}
				defer w.Close()

				_, err = w.Write([]byte(text))
				if err != nil {
					smtpu.log.WithError(err).Errorln("Message send failed")
					errChan <- err
					return
				}

				if err != nil {
					smtpu.log.WithError(err).Errorln("Message send failed")
					errChan <- err
					return
				}

				smtpu.log.WithField("id", messageID).Infoln("Message sent")
				smtpu.log.Infoln("Message sent")

				statusChan <- mttypes.SendStatus{
					RecipientID:  recipient.ID,
					TemplateName: templateName,
					Status:       "Sent",
				}

				errChan <- smtpu.msgStorage.PutValue(messageID, &mttypes.MessagesStorageValue{
					UserId:       recipient.ID,
					TemplateName: templateName,
					Variables:    recipient.Variables,
					CreatedAt:    time.Now().UTC(),
					Message:      base64.StdEncoding.EncodeToString([]byte(text)),
				})
			}(client, recipient, mailtext, msgID)
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
		smtpu.log.Info("Operation cancelled")
		err = ctx.Err()
	case <-mgDoneCh:
	}

	return resp, err
}

func (smtpu *smtpUpstream) SimpleSend(ctx context.Context, templateName string, tsv *mttypes.TemplateStorageValue, recipient *mttypes.Recipient) (status *mttypes.SendStatus, err error) {
	tmpl, err := smtpu.parseTemplate(templateName, tsv)
	if err != nil {
		return nil, err
	}

	text, err := smtpu.executeTemplate(tmpl, recipient, nil)
	if err != nil {
		return nil, err
	}

	tmplemail, err := texttemplate.New("emailtemplate").Parse(emailtemplate)
	if err != nil {
		return nil, err
	}

	msgID := time.Now().UTC().Format("20060102150405.") + strconv.Itoa(rand.Int())
	mailtext, err := smtpu.constructMessage(tmplemail, *recipient, msgID, tsv.Subject, text)
	if err != nil {
		return nil, err
	}

	mgDoneCh := make(chan struct{})
	go func() {
		defer close(mgDoneCh)

		host, _, _ := net.SplitHostPort(smtpu.smtpAddress)

		auth := smtp.PlainAuth("", smtpu.smtpLogin, smtpu.smtpPassword, host)

		tlsconfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         host,
		}

		conn, err := tls.Dial("tcp", smtpu.smtpAddress, tlsconfig)
		if err != nil {
			smtpu.log.WithError(err).Errorln("Message send failed")
			return
		}
		client, err := smtp.NewClient(conn, host)
		defer client.Quit()
		if err != nil {
			smtpu.log.WithError(err).Errorln("Message send failed")
			return
		}

		if err = client.Auth(auth); err != nil {
			smtpu.log.WithError(err).Errorln("Message send failed")
			return
		}
		if err = client.Mail(smtpu.senderMail); err != nil {
			smtpu.log.WithError(err).Errorln("Message send failed")
			return
		}
		if err = client.Rcpt(recipient.Email); err != nil {
			smtpu.log.WithError(err).Errorln("Message send failed")
			return
		}

		w, err := client.Data()
		if err != nil {
			smtpu.log.WithError(err).Errorln("Message send failed")
			return
		}
		defer w.Close()

		_, err = w.Write([]byte(mailtext))
		if err != nil {
			smtpu.log.WithError(err).Errorln("Message send failed")
			return
		}

		if err != nil {
			smtpu.log.WithError(err).Errorln("Message send failed")
			return
		}

		smtpu.log.WithField("id", msgID).Infoln("Message sent")

		status = &mttypes.SendStatus{
			RecipientID:  recipient.ID,
			TemplateName: templateName,
			Status:       "Sent",
		}

		err = smtpu.msgStorage.PutValue(msgID, &mttypes.MessagesStorageValue{
			UserId:       recipient.ID,
			TemplateName: templateName,
			Variables:    recipient.Variables,
			CreatedAt:    time.Now().UTC(),
			Message:      base64.StdEncoding.EncodeToString([]byte(text)),
		})
	}()

	select {
	case <-ctx.Done():
		smtpu.log.Info("Operation cancelled")
		err = ctx.Err()
	case <-mgDoneCh:
	}

	return status, err
}
