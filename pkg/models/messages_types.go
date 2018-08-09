package models

import "time"

// MessageGetResponse -- response to get message request
//
// swagger:model
type MessageGetResponse struct {
	ID string `json:"id"`
	// swagger: allOf
	*MessagesStorageValue
}

// MessagesStorageValue -- model for message in storage
//
// swagger:model
type MessagesStorageValue struct {
	UserID       string                 `json:"user_id"`
	TemplateName string                 `json:"template_name"`
	Variables    map[string]interface{} `json:"variables,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
	Message      string                 `json:"message"`
}

// MessageListEntry -- model for messages list
//
// swagger:model
type MessageListEntry struct {
	ID           string    `json:"id,omitempty"`
	UserID       string    `json:"user_id,omitempty"`
	TemplateName string    `json:"template_name,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
}

// MessageListEntry -- model for messages list with query
//
// swagger:model
type MessageListResponse struct {
	Messages []MessageListEntry `json:"message_list"`
	// swagger: allOf
	*MessageListQuery
}

// MessageListQuery -- query for message list
//
// swagger:model
type MessageListQuery struct {
	Page    int `form:"page,omitempty"`
	PerPage int `form:"per_page,omitempty"`
}
