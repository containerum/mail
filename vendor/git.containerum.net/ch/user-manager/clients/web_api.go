package clients

import (
	"net/http"

	"github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"
)

type WebAPIClient struct {
	log    *logrus.Entry
	client *resty.Client
}

type WebAPILoginRequest struct {
	Login    string `json:"login" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type WebAPIError struct {
	Message string `json:"message"`
}

func (e *WebAPIError) Error() string {
	return e.Message
}

func NewWebAPIClient(serverUrl string) *WebAPIClient {
	log := logrus.WithField("component", "web_api_client")
	client := resty.New().SetHostURL(serverUrl).SetLogger(log.WriterLevel(logrus.DebugLevel))
	return &WebAPIClient{
		log:    log,
		client: client,
	}
}

// returns raw answer from web-api
func (c *WebAPIClient) Login(request *WebAPILoginRequest) (ret map[string]interface{}, statusCode int, err error) {
	c.log.WithField("login", request.Login).Infoln("Signing in through web-api")
	ret = make(map[string]interface{})

	resp, err := c.client.R().SetBody(request).SetError(WebAPIError{}).SetResult(ret).Post("/api/login")
	if err != nil {
		c.log.WithError(err).Errorln("Sign in through web-api request failed")
		return nil, http.StatusInternalServerError, err
	}
	if resp.StatusCode() > 399 {
		msg := resp.Error().(*WebAPIError)
		c.log.WithField("message", msg.Message).Infoln("Sign in through web-api failed")
		return nil, resp.StatusCode(), msg
	}

	return ret, resp.StatusCode(), nil
}
