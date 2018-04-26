package clients

import (
	"context"

	"time"

	"git.containerum.net/ch/cherry"
	"git.containerum.net/ch/user-manager/pkg/models"
	"github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"
)

// UserManagerClient is interface to user-manager service
type UserManagerClient interface {
	UserInfoByID(ctx context.Context, userID string) (*models.User, error)
}

type httpUserManagerClient struct {
	log    *logrus.Entry
	client *resty.Client
}

// NewHTTPUserManagerClient returns rest-client to user-manager service
func NewHTTPUserManagerClient(serverURL string) UserManagerClient {
	log := logrus.WithField("component", "user_manager_client")
	client := resty.New().
		SetLogger(log.WriterLevel(logrus.DebugLevel)).
		SetHostURL(serverURL).
		SetDebug(true).
		SetTimeout(3 * time.Second).
		SetError(&cherry.Err{})
	client.JSONMarshal = jsoniter.Marshal
	client.JSONUnmarshal = jsoniter.Unmarshal
	return &httpUserManagerClient{
		log:    log,
		client: client,
	}
}

// NewHTTPUserManagerClient returns user info from user-manager
func (u *httpUserManagerClient) UserInfoByID(ctx context.Context, userID string) (*models.User, error) {
	u.log.WithField("id", userID).Info("Get user info from")
	ret := models.User{}

	resp, err := u.client.R().
		SetContext(ctx).
		SetResult(&ret).
		Get("/user/info/id/" + userID)
	if err != nil {
		return nil, err
	}

	if resp.Error() != nil {
		return nil, resp.Error().(*cherry.Err)
	}
	return &ret, nil
}
