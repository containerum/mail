package upstreams

import "bitbucket.org/exonch/ch-mail-templater/storages"

type Recipient struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Email     string            `json:"email"`
	Variables map[string]string `json:"variables"`
}

type SendRequest struct {
	Delay   int `json:"delay"` // in minutes
	Message struct {
		CommonVariables map[string]string `json:"common_variables"`
		Recipients      []Recipient       `json:"recipient_data"`
	} `json:"message"`
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
}
