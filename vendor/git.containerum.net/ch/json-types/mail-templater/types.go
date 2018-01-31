package mail

import "time"

type SimpleSendRequest struct {
	Template  string                 `json:"template" binding:"required"`
	UserID    string                 `json:"user_id" binding:"required,uuid4"`
	Variables map[string]interface{} `json:"variables"`
}

type SimpleSendResponse struct {
	UserID string `json:"user_id"`
}

type MessagesStorageValue struct {
	UserId       string                 `json:"user_id"`
	TemplateName string                 `json:"template_name"`
	Variables    map[string]interface{} `json:"variables,omitempty"`
	CreatedAt    time.Time              `json:"created_at"` // UTC
	Message      string                 `json:"message"`    // base64
}

type MessageGetResponse struct {
	Id string `json:"id"`
	*MessagesStorageValue
}

type TemplateStorageValue struct {
	Data      string    `json:"data"`
	Subject   string    `json:"template_subject"`
	CreatedAt time.Time `json:"created_at"` // UTC
}

type TemplateCreateRequest struct {
	Name    string `json:"template_name" binding:"required"`
	Version string `json:"template_version" binding:"required"`
	Data    string `json:"template_data" binding:"required,base64"`
	Subject string `json:"template_subject" binding:"required"`
}

type TemplateCreateResponse struct {
	Name    string `json:"template_name"`
	Version string `json:"template_version"`
}

type TemplateUpdateRequest struct {
	Data    string `json:"template_data" binding:"omitempty,base64"`
	Subject string `json:"template_subject" binding:"omitempty"`
}

type TemplateUpdateResponse struct {
	Name    string `json:"template_name"`
	Version string `json:"template_version"`
}

type TemplateDeleteResponse struct {
	Name    string `json:"template_name"`
	Version string `json:"template_version"`
}

type TemplatesDeleteResponse struct {
	Name string `json:"template_name"`
}

type Recipient struct {
	ID        string                 `json:"id" binding:"required,uuid4"`
	Name      string                 `json:"name" binding:"required"`
	Email     string                 `json:"email" binding:"required,email"`
	Variables map[string]interface{} `json:"variables"`
}

type SendRequestMessage struct {
	CommonVariables map[string]string `json:"common_variables"`
	Recipients      []Recipient       `json:"recipient_data" binding:"required"`
}

type SendRequest struct {
	Delay   int                `json:"delay" binding:"omitempty,min=0"` // in minutes
	Message SendRequestMessage `json:"message" binding:"required"`
}

type SendStatus struct {
	RecipientID  string `json:"recipient_id"`
	TemplateName string `json:"template_name"`
	Status       string `json:"status"`
}

type SendResponse struct {
	Statuses []SendStatus `json:"email_list"`
}

type TemplatesListEntry struct {
	Name     string   `json:"template_name"`
	Versions []string `json:"template_versions"`
}

type TemplatesListResponse struct {
	Templates []TemplatesListEntry `json:"templates"`
}

type MessageListQuery struct {
	Page    int `form:"page" binding:"required,gt=0"`
	PerPage int `form:"per_page" binding:"required,gt=0"`
}

type MessageListEntry struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	TemplateName string    `json:"template_name"`
	CreatedAt    time.Time `json:"created_at"`
}

type MessageListResponse struct {
	Messages []MessageListEntry `json:"message_list"`
}