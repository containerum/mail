package clients

import (
	"fmt"
	"strconv"

	"github.com/google/go-github/github"
	"github.com/huandu/facebook"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	githubOAuth "golang.org/x/oauth2/github"
	googleOAuth "golang.org/x/oauth2/google"
	google "google.golang.org/api/oauth2/v2"
)

type OAuthUserInfo struct {
	UserID string
	Email  string
}

type OAuthClient interface {
	GetUserInfo(accessToken string) (info *OAuthUserInfo, err error)
	GetResource() OAuthResource
}

type OAuthResource string

const (
	GitHubOAuth   OAuthResource = "github"
	GoogleOAuth   OAuthResource = "google"
	FacebookOAuth OAuthResource = "facebook"
)

var oAuthClients = make(map[OAuthResource]OAuthClient)

func OAuthClientByResource(resource OAuthResource) (client OAuthClient, exists bool) {
	client, exists = oAuthClients[resource]
	return
}

func RegisterOAuthClient(client OAuthClient) {
	oAuthClients[client.GetResource()] = client
}

type GithubOAuthClient struct {
	log         *logrus.Entry
	oAuthConfig *oauth2.Config
}

func NewGithubOAuthClient(appID, appSecret string) *GithubOAuthClient {
	return &GithubOAuthClient{
		log: logrus.WithField("component", "github_oauth"),
		oAuthConfig: &oauth2.Config{
			ClientID:     appID,
			ClientSecret: appSecret,
			Endpoint:     githubOAuth.Endpoint,
			Scopes:       []string{string(github.ScopeUser), string(github.ScopeUserEmail)},
		},
	}
}

func (gh *GithubOAuthClient) GetResource() OAuthResource {
	return GitHubOAuth
}

func (gh *GithubOAuthClient) GetUserInfo(accessToken string) (info *OAuthUserInfo, err error) {
	gh.log.WithField("token", accessToken).Infoln("Get GitHub user info")
	ctx := context.Background()
	ts := gh.oAuthConfig.TokenSource(ctx, &oauth2.Token{AccessToken: accessToken})
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	user, resp, err := client.Users.Get(ctx, "")
	if err != nil {
		gh.log.WithError(err).Errorln("Request error")
		return nil, err
	}
	if resp.StatusCode >= 400 {
		gh.log.WithField("error", resp.Status).Errorf("GitHub API error")
		return nil, fmt.Errorf("github API error")
	}

	return &OAuthUserInfo{
		UserID: strconv.Itoa(user.GetID()),
		Email:  user.GetEmail(),
	}, nil
}

type GoogleOAuthClient struct {
	log         *logrus.Entry
	oAuthConfig *oauth2.Config
}

func NewGoogleOAuthClient(appID, appSecret string) *GoogleOAuthClient {
	return &GoogleOAuthClient{
		log: logrus.WithField("component", "google_oauth"),
		oAuthConfig: &oauth2.Config{
			ClientID:     appID,
			ClientSecret: appSecret,
			Endpoint:     googleOAuth.Endpoint,
			Scopes:       []string{google.UserinfoProfileScope, google.UserinfoEmailScope},
		},
	}
}

func (gc *GoogleOAuthClient) GetResource() OAuthResource {
	return GoogleOAuth
}

func (gc *GoogleOAuthClient) GetUserInfo(accessToken string) (info *OAuthUserInfo, err error) {
	gc.log.WithField("token", accessToken).Infoln("Get Google user info")
	ctx := context.Background()
	ts := gc.oAuthConfig.TokenSource(ctx, &oauth2.Token{AccessToken: accessToken})
	tc := oauth2.NewClient(ctx, ts)

	client, err := google.New(tc)
	if err != nil {
		gc.log.WithError(err).Errorln("Client create failed")
		return nil, err
	}

	googleInfo, err := google.NewUserinfoV2MeService(client).Get().Do()
	if err != nil {
		gc.log.WithError(err).Errorln("Fetch user info failed")
		return nil, err
	}

	return &OAuthUserInfo{
		UserID: googleInfo.Id,
		Email:  googleInfo.Email,
	}, nil
}

type FacebookOAuthClient struct {
	log *logrus.Entry
	app *facebook.App
}

func NewFacebookOAuthClient(appID, appSecret string) *FacebookOAuthClient {
	return &FacebookOAuthClient{
		log: logrus.WithField("component", "facebook_oauth"),
		app: facebook.New(appID, appSecret),
	}
}

func (fb *FacebookOAuthClient) GetResource() OAuthResource {
	return FacebookOAuth
}

func (fb *FacebookOAuthClient) GetUserInfo(accessToken string) (info *OAuthUserInfo, err error) {
	fb.log.WithField("token", accessToken).Infoln("Get Facebook user info")

	session := fb.app.Session(accessToken)

	resp, err := session.Get("/me", facebook.Params{
		"access_token": accessToken,
		"fields":       "id,email",
	})
	if err != nil {
		fb.log.WithError(err).Errorln("Fetch user info failed")
		return nil, err
	}

	return &OAuthUserInfo{
		UserID: resp.Get("id").(string),
		Email:  resp.Get("email").(string),
	}, nil
}
