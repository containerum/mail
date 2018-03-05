package auth

type ExtendTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type CheckTokenResponse struct {
	Access struct {
		Namespace []Resource `json:"namespace"`
		Volume    []Resource `json:"volume"`
	} `json:"access"`
}
