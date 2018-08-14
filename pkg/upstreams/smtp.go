package upstreams

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"html/template"
	"strings"
	texttemplate "text/template"

	"time"

	"crypto/tls"
	"net/smtp"
	"strconv"

	"net"

	"git.containerum.net/ch/mail-templater/pkg/models"
	"git.containerum.net/ch/mail-templater/pkg/storages"
	"git.containerum.net/ch/mail-templater/pkg/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
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

// NewSMTPUpstream returns mail upstream which uses SMTP service to send emails.
func NewSMTPUpstream(msgStorage storages.MessagesStorage, senderName string, senderMail string, smtpAddress string, smtpLogin string, smtpPassword string) Upstream {
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

func (smtpu *smtpUpstream) executeTemplate(tmpl *template.Template, recipient *models.Recipient, commonVars map[string]string) (string, error) {
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

func (smtpu *smtpUpstream) constructMessage(template *texttemplate.Template, recipient models.Recipient, msgID string, subject string, text string) (*string, error) {
	newmail := mailData{SenderName: smtpu.senderName,
		SenderMail:    smtpu.senderMail,
		RecipientName: recipient.Name,
		RecipientMail: recipient.Email,
		Subject:       subject,
		MessageID:     msgID,
		Body:          text}

	var mailtext bytes.Buffer
	if err := template.Execute(&mailtext, newmail); err != nil {
		return nil, err
	}

	msg := mailtext.String()
	return &msg, nil
}

func (smtpu *smtpUpstream) parseTemplate(templateName string, tsv *models.Template) (tmpl *template.Template, err error) {
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

func (smtpu *smtpUpstream) newSMTPClient(recipientEmail string, text string) error {
	host, _, _ := net.SplitHostPort(smtpu.smtpAddress)

	conn, err := tls.Dial("tcp", smtpu.smtpAddress, &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         host,
	})
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	defer client.Quit()

	if err = client.Auth(smtp.PlainAuth("", smtpu.smtpLogin, smtpu.smtpPassword, host)); err != nil {
		return err
	}
	if err = client.Mail(smtpu.senderMail); err != nil {
		return err
	}

	if err = client.Rcpt(recipientEmail); err != nil {
		return err
	}

	w, err := client.Data()
	if err != nil {
		return err
	}
	defer w.Close()

	if _, err := w.Write([]byte(text)); err != nil {
		return err
	}
	return nil
}

//Send
//Sends email using smtp
func (smtpu *smtpUpstream) Send(ctx context.Context, templateName string, tsv *models.Template, request *models.SendRequest) (*models.SendResponse, error) {
	var errs []string
	var msgNumber = 0

	resp := &models.SendResponse{Statuses: make([]models.SendStatus, 0)}

	tmpl, err := smtpu.parseTemplate(templateName, tsv)
	if err != nil {
		return nil, err
	}

	tmplemail, err := texttemplate.New("emailtemplate").Parse(emailtemplate)
	if err != nil {
		return nil, err
	}

	var g errgroup.Group
	for _, r := range request.Message.Recipients {
		msgNumber++

		recipient := r
		var text string
		text, err = smtpu.executeTemplate(tmpl, &recipient, request.Message.CommonVariables)
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}

		messageID := time.Now().UTC().Format("20060102150405.123456.") + strconv.Itoa(msgNumber)
		mailtext, err := smtpu.constructMessage(tmplemail, recipient, messageID, tsv.Subject, text)
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}

		g.Go(func() error {
			if err := smtpu.newSMTPClient(recipient.Email, *mailtext); err != nil {
				smtpu.log.WithError(err).Errorln("Message send failed")
				return err
			}

			resp.Statuses = append(resp.Statuses, models.SendStatus{
				RecipientID:  recipient.ID,
				TemplateName: templateName,
				Status:       "Sent",
			})
			smtpu.log.WithField("id", messageID).Infoln("Message sent")

			if err := smtpu.msgStorage.PutMessage(messageID, &models.MessagesStorageValue{
				UserID:       recipient.ID,
				TemplateName: templateName,
				Variables:    recipient.Variables,
				CreatedAt:    time.Now().UTC(),
				Message:      base64.StdEncoding.EncodeToString([]byte(text)),
			}); err != nil {
				return err
			}
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}

	if len(errs) > 0 {
		err = errors.New(strings.Join(errs, "; "))
	}
	return resp, err
}

//SimpleSend
//Sends email using smtp in simple way
func (smtpu *smtpUpstream) SimpleSend(ctx context.Context, templateName string, tsv *models.Template, recipient *models.Recipient) (*models.SendStatus, error) {
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

	messageID := time.Now().UTC().Format("20060102150405.123456.") + "1"
	mailtext, err := smtpu.constructMessage(tmplemail, *recipient, messageID, tsv.Subject, text)
	if err != nil {
		return nil, err
	}

	if err = smtpu.newSMTPClient(recipient.Email, *mailtext); err != nil {
		smtpu.log.WithError(err).Errorln("Message send failed")
		return nil, err
	}

	if err := smtpu.msgStorage.PutMessage(messageID, &models.MessagesStorageValue{
		UserID:       recipient.ID,
		TemplateName: templateName,
		Variables:    recipient.Variables,
		CreatedAt:    time.Now().UTC(),
		Message:      base64.StdEncoding.EncodeToString([]byte(text)),
	}); err != nil {
		return nil, err
	}

	return &models.SendStatus{
		RecipientID:  recipient.ID,
		TemplateName: templateName,
		Status:       "Sent",
	}, err
}

func (smtpu *smtpUpstream) CheckStatus() error {
	smtpu.log.Debugln("Checking SMTP server connection")
	return utils.Retry(3, 15*time.Second, func() error {
		host, _, _ := net.SplitHostPort(smtpu.smtpAddress)

		conn, err := tls.Dial("tcp", smtpu.smtpAddress, &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         host,
		})
		if err != nil {
			operr, ok := err.(*net.OpError)
			if ok {
				if operr.Temporary() || operr.Timeout() {
					return err
				}
			}
			return &utils.StopRetry{Err: err}
		}
		defer conn.Close()

		client, err := smtp.NewClient(conn, host)
		if err != nil {
			return &utils.StopRetry{Err: err}
		}
		defer client.Quit()

		if err = client.Auth(smtp.PlainAuth("", smtpu.smtpLogin, smtpu.smtpPassword, host)); err != nil {
			return &utils.StopRetry{Err: err}
		}
		return nil
	})
}
