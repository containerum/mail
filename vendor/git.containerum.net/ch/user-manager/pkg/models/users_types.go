package models

// RegisterRequest -- request to create new user
//
// swagger:model
type RegisterRequest struct {
	// required: true
	Login string `json:"login"`
	// required: true
	Password string `json:"password"`
	// required: true
	ReCaptcha string `json:"recaptcha"`
	Referral  string `json:"referral"`
}

// UserList -- users list
//
// swagger:model
type UserList struct {
	Users []User `json:"users,omitempty"`
	Pages uint   `json:"pages,omitempty"`
}

// User -- user model
//
// swagger:model
type User struct {
	// swagger: allOf
	*UserLogin
	// swagger: allOf
	*Accounts
	// swagger: allOf
	*Profile
	Role          string `json:"role,omitempty"`
	IsActive      bool   `json:"is_active,omitempty"`
	IsInBlacklist bool   `json:"is_in_blacklist,omitempty"`
	IsDeleted     bool   `json:"is_deleted,omitempty"`
}

// UserList -- model for user login, password and id
//
// swagger:model
type UserLogin struct {
	ID       string `json:"id,omitempty"`
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}

// Accounts -- list of bound accounts
//
// swagger:model
type Accounts struct {
	Accounts map[string]string `json:"accounts,omitempty"`
}

// Profile -- additional user information
//
// swagger:model
type Profile struct {
	Referral      string   `json:"referral,omitempty"`
	Access        string   `json:"access,omitempty"`
	CreatedAt     string   `json:"created_at,omitempty"`
	DeletedAt     string   `json:"deleted_at,omitempty"`
	BlacklistedAt string   `json:"blacklisted_at,omitempty"`
	LastLogin     string   `json:"last_login,omitempty"`
	Data          UserData `json:"data,omitempty"`
}

// BoundAccountDeleteRequest -- request to remove bound account
//
// swagger:model
type BoundAccountDeleteRequest struct {
	// required: true
	Resource string `json:"resource"`
}

// BoundAccounts -- bound accounts list for user
//
// swagger:model
type BoundAccounts map[string]string

// UserData -- user profile data
//
// swagger:model
type UserData map[string]interface{}

// LoginID -- logins and user ID
//
// swagger:model
type LoginID map[string]string

// IDList -- ids list
//
// swagger:model
type IDList []string
