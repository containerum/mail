package upstreams

import (
	"context"

	mttypes "git.containerum.net/ch/json-types/mail-templater"
)

type Upstream interface {
	Send(ctx context.Context, templateName string, tsv *mttypes.TemplateStorageValue, request *mttypes.SendRequest) (resp *mttypes.SendResponse, err error)
	SimpleSend(ctx context.Context, templateName string, tsv *mttypes.TemplateStorageValue, recipient *mttypes.Recipient) (status *mttypes.SendStatus, err error)
}
