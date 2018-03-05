package user

import "time"

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

type UserInfoByIDGetResponse struct {
	Login string                 `json:"login"`
	Role  string                 `json:"role"`
	Data  map[string]interface{} `json:"data"`
}

type UserInfoByLoginGetResponse struct {
	ID string                 	 `json:"id"`
	Role  string                 `json:"role"`
	Data  map[string]interface{} `json:"data"`
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

type UserListQuery struct {
	Page    int `form:"page" binding:"required,gt=0"`
	PerPage int `form:"per_page" binding:"required,gt=0"`
}

type CompleteDeleteHandlerRequest struct {
	UserID string `json:"user_id" binding:"required,uuid4"`
}
