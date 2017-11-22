package upstreams

type Recipient struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Email     string            `json:"email"`
	Variables map[string]string `json:"variables"`
}

type SendRequest struct {
	Delay   int `json:"delay"` // in minutes
	Message struct {
		Subject         string            `json:"subject"`
		SenderEmail     string            `json:"sender_email"`
		SenderName      string            `json:"sender_name"`
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
	Send(templateName, templateContent string, request *SendRequest) (resp *SendResponse, err error)
}
