package mail

//Types related to sending emails

type SimpleSendRequest struct {
	Template  string                 `json:"template"`
	UserID    string                 `json:"user_id"`
	Variables map[string]interface{} `json:"variables"`
}

type SimpleSendResponse struct {
	UserID string `json:"user_id"`
}

type Recipient struct {
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	Email     string                 `json:"email"`
	Variables map[string]interface{} `json:"variables"`
}

type SendRequest struct {
	Delay   int `json:"delay"` // in minutes
	Message struct {
		CommonVariables map[string]string `json:"common_variables"`
		Recipients      []Recipient       `json:"recipient_data"`
	} `json:"message"`
}

type SendResponse struct {
	Statuses []SendStatus `json:"email_list"`
}

type SendStatus struct {
	RecipientID  string `json:"recipient_id"`
	TemplateName string `json:"template_name"`
	Status       string `json:"status"`
}
