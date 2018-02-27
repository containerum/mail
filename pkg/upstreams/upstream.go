package upstreams

import (
	"context"

	mttypes "git.containerum.net/ch/json-types/mail-templater"
)

// Upstream is interface for sending email
type Upstream interface {
	// Send sends multiple emails to multiple clients
	Send(ctx context.Context, templateName string, tsv *mttypes.TemplateStorageValue, request *mttypes.SendRequest) (resp *mttypes.SendResponse, err error)

	// SimpleSend sends email to one client
	SimpleSend(ctx context.Context, templateName string, tsv *mttypes.TemplateStorageValue, recipient *mttypes.Recipient) (status *mttypes.SendStatus, err error)
}
