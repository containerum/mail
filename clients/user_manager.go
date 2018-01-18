package clients

import (
	"context"

	"git.containerum.net/ch/json-types/errors"
	umtypes "git.containerum.net/ch/json-types/user-manager"
	"github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"
)

// UserManagerClient is interface to user-manager service
type UserManagerClient interface {
	UserInfoByID(ctx context.Context, userID string) (*umtypes.UserInfoGetResponse, error)
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
		SetError(errors.Error{})
	client.JSONMarshal = jsoniter.Marshal
	client.JSONUnmarshal = jsoniter.Unmarshal
	return &httpUserManagerClient{
		log:    log,
		client: client,
	}
}

func (u *httpUserManagerClient) UserInfoByID(ctx context.Context, userID string) (*umtypes.UserInfoGetResponse, error) {
	u.log.WithField("id", userID).Info("Get user info from")
	ret := umtypes.UserInfoGetResponse{}
	resp, err := u.client.R().
		SetContext(ctx).
		SetHeader(umtypes.UserIDHeader, userID).
		SetResult(&ret).
		Get("/user/info")
	if err != nil {
		return nil, err
	}
	if err := resp.Error(); err != nil {
		return nil, err.(*errors.Error)
	}
	return &ret, nil
}
