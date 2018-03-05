package mail

import "time"

//Types related to storing messages

type MessageGetResponse struct {
	Id string `json:"id"`
	*MessagesStorageValue
}

type MessagesStorageValue struct {
	UserId       string                 `json:"user_id"`
	TemplateName string                 `json:"template_name"`
	Variables    map[string]interface{} `json:"variables,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
	Message      string                 `json:"message"`
}

type MessageListEntry struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	TemplateName string    `json:"template_name"`
	CreatedAt    time.Time `json:"created_at"`
}

type MessageListResponse struct {
	Messages []MessageListEntry `json:"message_list"`
	*MessageListQuery
}

type MessageListQuery struct {
	Page    int `form:"page,omitempty"`
	PerPage int `form:"per_page,omitempty"`
}
