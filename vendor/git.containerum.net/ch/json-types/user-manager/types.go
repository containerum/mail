package user

import (
	"time"
)

type ProfileData struct {
	Email          string `json:"email,omitempty" binding:"omitempty,email"`
	Address        string `json:"address,omitempty"`
	Phone          string `json:"phone,omitempty"`
	FirstName      string `json:"first_name,omitempty"`
	LastName       string `json:"last_name,omitempty"`
	IsOrganization bool   `json:"is_organization,omitempty"`
	TaxCode        string `json:"tax_code,omitempty"`
	Company        string `json:"company,omitempty"`
}

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

type UserRole int

const (
	RoleUser UserRole = iota
	RoleAdmin
)

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

type ActivateRequest struct {
	Link string `json:"link" binding:"required"`
}

type ResendLinkRequest struct {
	UserName string `json:"username" binding:"required,email"`
}

type UserInfoByIDGetResponse struct {
	Login string      `json:"login"`
	Data  ProfileData `json:"data"`
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
	Login     string      `json:"login"`
	Data      ProfileData `json:"data"`
	ID        string      `json:"id"`
	IsActive  bool        `json:"is_active"`
	CreatedAt time.Time   `json:"created_at"`
}

type UserListEntry struct {
	ID            string      `json:"id"`
	Login         string      `json:"login"`
	Referral      string      `json:"referral"`
	Role          UserRole    `json:"role"`
	Access        string      `json:"access"`
	CreatedAt     time.Time   `json:"created_at"`
	DeletedAt     time.Time   `json:"deleted_at"`
	BlacklistedAt time.Time   `json:"blacklisted_at"`
	Data          ProfileData `json:"data"`
	IsActive      bool        `json:"is_active"`
	IsInBlacklist bool        `json:"is_in_blacklist"`
	IsDeleted     bool        `json:"is_deleted"`
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
