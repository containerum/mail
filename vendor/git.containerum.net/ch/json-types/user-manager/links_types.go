package user

import "time"

type LinkType string

const (
	LinkTypeConfirm   LinkType = "confirm"
	LinkTypePwdChange LinkType = "pwd_change"
	LinkTypeDelete    LinkType = "delete"
)

type Link struct {
	Link      string    `json:"link"`
	Type      LinkType  `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	ExpiredAt time.Time `json:"expired_at"`
	IsActive  bool      `json:"is_active"`
	SentAt    time.Time `json:"sent_at,omitempty"`
}

type LinksGetResponse struct {
	Links []Link `json:"links"`
}

type ResendLinkRequest struct {
	UserName string `json:"username" binding:"required,email"`
}

type ActivateRequest struct {
	Link string `json:"link" binding:"required"`
}