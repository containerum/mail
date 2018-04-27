package models

// SimpleSendRequest -- request to send mail using simple send method
//
// swagger:model
type SimpleSendRequest struct {
	// required: true
	Template string `json:"template"`
	// required: true
	UserID string `json:"user_id"`
	// required: true
	Variables map[string]interface{} `json:"variables"`
}

// SimpleSendResponse -- responce to send mail using simple send method
//
// swagger:model
type SimpleSendResponse struct {
	UserID string `json:"user_id"`
}

// Recipient -- recipient info for send mail method
//
// swagger:model
type Recipient struct {
	// required: true
	ID string `json:"id"`
	// required: true
	Name string `json:"name"`
	// required: true
	Email     string                 `json:"email"`
	Variables map[string]interface{} `json:"variables"`
}

// SimpleSendResponse -- responce to send mail using send method
//
// swagger:model
type SendRequest struct {
	Delay   int `json:"delay"` // in minutes
	Message struct {
		CommonVariables map[string]string `json:"common_variables"`
		// required: true
		Recipients []Recipient `json:"recipient_data"`
	} `json:"message"`
}

// SendResponse -- responce to send mail using send method
//
// swagger:model
type SendResponse struct {
	Statuses []SendStatus `json:"email_list"`
}

// SendStatus -- status of sent emails
//
// swagger:model
type SendStatus struct {
	RecipientID  string `json:"recipient_id"`
	TemplateName string `json:"template_name"`
	Status       string `json:"status"`
}
