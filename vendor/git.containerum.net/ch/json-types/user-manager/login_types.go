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