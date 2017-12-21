package clients

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"
	um "git.containerum.net/ch/user-manager/routes"
	chutils "git.containerum.net/ch/utils"
)

type UserManagerClient struct {
	log *logrus.Entry
	client *resty.Client
}

func NewUserManagerClient(serverUrl string) *UserManagerClient {
	log := logrus.WithField("component", "user_manager_client")
	client := resty.New().SetLogger(log.WriterLevel(logrus.DebugLevel)).SetHostURL(serverUrl)
	return &UserManagerClient{
		log: log,
		client: client,
	}
}

func (u *UserManagerClient) UserInfoByID(userID string) (*um.UserInfoGetResponse, error) {
	u.log.WithField("id", userID).Info("Get user info from")
	ret := um.UserInfoGetResponse{}
	resp, err := u.client.R().
		SetHeader(um.UserIDHeader, userID).
		SetResult(&ret).
		SetError(chutils.Error{}).
		Get("/user/info")
	if err != nil {
		return nil, err
	}
	return &ret, resp.Error().(*chutils.Error)
}