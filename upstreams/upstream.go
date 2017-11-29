package upstreams

import "git.containerum.net/ch/mail-templater/storages"

type Recipient struct {
	ID        string            `json:"id" binding:"required,uuid4"`
	Name      string            `json:"name" binding:"required"`
	Email     string            `json:"email" binding:"required,email"`
	Variables map[string]string `json:"variables" binding:"required"`
}

type SendRequest struct {
	Delay   int `json:"delay" binding:"required,min=0"` // in minutes
	Message struct {
		CommonVariables map[string]string `json:"common_variables" binding:"required"`
		Recipients      []Recipient       `json:"recipient_data" binding:"required"`
	} `json:"message" binding:"required"`
}

type SendStatus struct {
	RecipientID  string `json:"recipient_id"`
	TemplateName string `json:"template_name"`
	Status       string `json:"status"`
}

type SendResponse struct {
	Statuses []SendStatus `json:"email_list"`
}

type Upstream interface {
	Send(templateName string, tsv *storages.TemplateStorageValue, request *SendRequest) (resp *SendResponse, err error)
	SimpleSend(templateName string, tsv *storages.TemplateStorageValue, recipient *Recipient) (status *SendStatus, err error)
}
