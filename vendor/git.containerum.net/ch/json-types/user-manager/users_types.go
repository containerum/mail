package user

type RegisterRequest struct {
	Login     string `json:"login"`
	Password  string `json:"password"`
	ReCaptcha string `json:"recaptcha"`
	Referral  string `json:"referral"`
}

type UserList struct {
	Users []User `json:"users,omitempty"`
}

type User struct {
	*UserLogin
	*Accounts
	*Profile
	Role          string `json:"role,omitempty"`
	IsActive      bool   `json:"is_active,omitempty"`
	IsInBlacklist bool   `json:"is_in_blacklist,omitempty"`
	IsDeleted     bool   `json:"is_deleted,omitempty"`
}

type UserLogin struct {
	ID    string `json:"id,omitempty"`
	Login string `json:"login,omitempty"`
}

type Accounts struct {
	Accounts map[string]string `json:"accounts,omitempty"`
}

type Profile struct {
	Referral      string                 `json:"referral,omitempty"`
	Access        string                 `json:"access,omitempty"`
	CreatedAt     string                 `json:"created_at,omitempty"`
	DeletedAt     string                 `json:"deleted_at,omitempty"`
	BlacklistedAt string                 `json:"blacklisted_at,omitempty"`
	Data          map[string]interface{} `json:"data,omitempty"`
}

type BoundAccountDeleteRequest struct {
	Resource string `json:"resource"`
}