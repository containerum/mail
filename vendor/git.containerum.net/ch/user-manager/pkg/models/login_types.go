package models

// LoginRequest -- login request (for basic login)
//
// swagger:model
type LoginRequest struct {
	// required: true
	Login string `json:"login"`
	// required: true
	Password string `json:"password"`
}

// LoginRequest -- login request (for token login)
//
// swagger:model
type OneTimeTokenLoginRequest struct {
	// required: true
	Token string `json:"token"`
}

// LoginRequest -- login request (for oauth login)
//
// swagger:model
type OAuthLoginRequest struct {
	// required: true
	Resource OAuthResource `json:"resource"`
	// required: true
	AccessToken string `json:"access_token"`
}
