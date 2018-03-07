package user

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type OneTimeTokenLoginRequest struct {
	Token string `json:"token"`
}

type OAuthLoginRequest struct {
	Resource    OAuthResource `json:"resource"`
	AccessToken string        `json:"access_token"`
}

type WebAPILoginResponse struct {
	BadSignature bool   `json:"bad_signature"`
	IsExpired    bool   `json:"is_expired"`
	Token        string `json:"token"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         struct {
		CreatedAt string                 `json:"created_at"`
		Data      map[string]interface{} `json:"data"`
		ID        string                 `json:"id"`
		IsActive  bool                   `json:"is_active"`
		Login     string                 `json:"login"`
	} `json:"user"`
}

type WebAPIResource struct {
	ID     string `json:"id"`
	Name   string `json:"label"`
	Access string `json:"access"`
}
