package user

import "time"

type LinkType string

const (
	LinkTypeConfirm   LinkType = "confirm"
	LinkTypePwdChange LinkType = "pwd_change"
	LinkTypeDelete    LinkType = "delete"
)

type Link struct {
	Link      string    `json:"link,omitempty"`
	Type      LinkType  `json:"type,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	ExpiredAt time.Time `json:"expired_at,omitempty"`
	IsActive  bool      `json:"is_active,omitempty"`
	SentAt    time.Time `json:"sent_at,omitempty"`
}

type Links struct {
	Links []Link `json:"links,omitempty"`
}