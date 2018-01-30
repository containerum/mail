package user

import (
	"time"
)

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

type BasicLoginRequest struct {
	Username  string `json:"username" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
	ReCaptcha string `json:"recaptcha" binding:"required"`
}

type OneTimeTokenLoginRequest struct {
	Token string `json:"token" binding:"required"`
}

type OAuthLoginRequest struct {
	Resource    OAuthResource `json:"resource" binding:"required"`
	AccessToken string        `json:"access_token" binding:"required"`
}

type PasswordChangeRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
}

type PasswordRestoreRequest struct {
	Link        string `json:"link" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type PasswordResetRequest struct {
	Username string `json:"username" binding:"required,email"`
}

type UserCreateRequest struct {
	UserName  string `json:"username" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
	Referral  string `json:"referral" binding:"omitempty,url"`
	ReCaptcha string `json:"recaptcha" binding:"required"`
}

type UserCreateResponse struct {
	ID       string `json:"id"`
	Login    string `json:"login"`
	IsActive bool   `json:"is_active"`
}

type UserCreateWebAPIRequest struct {
	ID        string                 `json:"id"`
	UserName  string                 `json:"username" binding:"required,email"`
	Password  string                 `json:"password" binding:"required"`
	Data      map[string]interface{} `json:"data"`
	CreatedAt string                 `json:"created_at"`
	IsActive  bool                   `json:"is_active"`
}

type ActivateRequest struct {
	Link string `json:"link" binding:"required"`
}

type ResendLinkRequest struct {
	UserName string `json:"username" binding:"required,email"`
}

type UserInfoByIDGetResponse struct {
	Login string                 `json:"login"`
	Data  map[string]interface{} `json:"data"`
}

type BlacklistedUserEntry struct {
	Login string `json:"login"`
	ID    string `json:"id"`
}

type BlacklistGetResponse struct {
	BlacklistedUsers []BlacklistedUserEntry `json:"blacklist_users"`
}

type LinksGetResponse struct {
	Links []Link `json:"links"`
}

type UserInfoGetResponse struct {
	Login     string                 `json:"login"`
	Data      map[string]interface{} `json:"data"`
	ID        string                 `json:"id"`
	Role      string                 `json:"role"`
	IsActive  bool                   `json:"is_active"`
	CreatedAt time.Time              `json:"created_at"`
}

type UserListEntry struct {
	ID            string                 `json:"id"`
	Login         string                 `json:"login"`
	Referral      string                 `json:"referral"`
	Role          string                 `json:"role"`
	Access        string                 `json:"access"`
	CreatedAt     string                 `json:"created_at"`
	DeletedAt     string                 `json:"deleted_at"`
	BlacklistedAt string                 `json:"blacklisted_at"`
	Data          map[string]interface{} `json:"data"`
	IsActive      bool                   `json:"is_active"`
	IsInBlacklist bool                   `json:"is_in_blacklist"`
	IsDeleted     bool                   `json:"is_deleted"`
	Accounts      map[string]string      `json:"accounts"`
}

type UserListGetResponse struct {
	Users []UserListEntry `json:"users"`
}

type WebAPILoginRequest struct {
	Username string `json:"username" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserToBlacklistRequest struct {
	UserID string `json:"user_id" binding:"required,uuid4"`
}

type UserListQuery struct {
	Page    int `form:"page" binding:"required,gt=0"`
	PerPage int `form:"per_page" binding:"required,gt=0"`
}

type CompleteDeleteHandlerRequest struct {
	UserID string `json:"user_id" binding:"required,uuid4"`
}

type BoundAccountsResponce struct {
	Accounts map[string]string `json:"accounts" binding:"required"`
}

type BoundAccountDeleteRequest struct {
	Resource string `json:"resource" binding:"required"`
}

type DomainToBlacklistRequest struct {
	Domain string `json:"domain" binding:"required"`
}

type DomainListResponce struct {
	DomainList []DomainResponce `json:"domain_list"`
}

type DomainResponce struct {
	Domain    string `json:"domain"`
	AddedBy   string `json:"added_by"`
	CreatedAt string `json:"created_at"`
}
