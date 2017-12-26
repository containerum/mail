package clients

import (
	"git.containerum.net/ch/json-types/errors"
	umtypes "git.containerum.net/ch/json-types/user-manager"
	"github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"
)

type UserManagerClient struct {
	log    *logrus.Entry
	client *resty.Client
}

func NewUserManagerClient(serverUrl string) *UserManagerClient {
	log := logrus.WithField("component", "user_manager_client")
	client := resty.New().SetLogger(log.WriterLevel(logrus.DebugLevel)).SetHostURL(serverUrl)
	return &UserManagerClient{
		log:    log,
		client: client,
	}
}

func (u *UserManagerClient) UserInfoByID(userID string) (*umtypes.UserInfoGetResponse, error) {
	u.log.WithField("id", userID).Info("Get user info from")
	ret := umtypes.UserInfoGetResponse{}
	resp, err := u.client.R().
		SetHeader(umtypes.UserIDHeader, userID).
		SetResult(&ret).
		SetError(errors.Error{}).
		Get("/user/info")
	if err != nil {
		return nil, err
	}
	return &ret, resp.Error().(*errors.Error)
}
