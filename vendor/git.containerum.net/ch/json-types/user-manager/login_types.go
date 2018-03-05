package user

type BasicLoginRequest struct {
	Username  string `json:"username" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
}

type OneTimeTokenLoginRequest struct {
	Token string `json:"token" binding:"required"`
}

type OAuthLoginRequest struct {
	Resource    OAuthResource `json:"resource" binding:"required"`
	AccessToken string        `json:"access_token" binding:"required"`
}

type WebAPILoginRequest struct {
	Username string `json:"username" binding:"required,email"`
	Password string `json:"password" binding:"required"`
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

type WebAPIVolumesResponse struct {
	Name string `json:"name"`
}

type WebAPINamespaceResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type WebAPIResource struct {
	ID     string `json:"id"`
	Label  string `json:"label"`
	Access string `json:"access"`
}
