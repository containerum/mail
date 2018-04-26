package models

type OAuthResource string

const (
	GitHubOAuth   OAuthResource = "github"
	GoogleOAuth   OAuthResource = "google"
	FacebookOAuth OAuthResource = "facebook"
)
