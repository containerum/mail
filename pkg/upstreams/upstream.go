package upstreams

import (
	"context"

	"git.containerum.net/ch/mail-templater/pkg/models"
)

// Upstream is interface for sending email
type Upstream interface {
	// Send sends multiple emails to multiple clients
	Send(ctx context.Context, templateName string, tsv *models.Template, request *models.SendRequest) (resp *models.SendResponse, err error)

	// SimpleSend sends email to one client
	SimpleSend(ctx context.Context, templateName string, tsv *models.Template, recipient *models.Recipient) (status *models.SendStatus, err error)

	CheckStatus() (err error)
}
