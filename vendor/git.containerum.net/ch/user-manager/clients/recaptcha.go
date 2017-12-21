package clients

import (
	"time"

	"net/url"

	"github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"
)

const reCaptchaAPI = "https://www.google.com/recaptcha/api"

type ReCaptchaClient struct {
	client     *resty.Client
	log        *logrus.Entry
	privateKey string
}

type ReCaptchaResponse struct {
	Success     bool      `json:"success"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []int     `json:"error-codes"`
}

func NewReCaptchaClient(privateKey string) *ReCaptchaClient {
	log := logrus.WithField("component", "recaptcha")
	client := resty.New().SetLogger(log.WriterLevel(logrus.DebugLevel)).SetHostURL(reCaptchaAPI)
	return &ReCaptchaClient{
		log:        log,
		client:     client,
		privateKey: privateKey,
	}
}

func (c *ReCaptchaClient) Check(remoteIP, clientResponse string) (r *ReCaptchaResponse, err error) {
	c.log.Infoln("Checking ReCaptcha from", remoteIP)
	r = new(ReCaptchaResponse)
	_, err = c.client.R().SetResult(r).SetMultiValueFormData(url.Values{
		"secret":   {c.privateKey},
		"remoteip": {remoteIP},
		"response": {clientResponse},
	}).Post("/siteverify")
	return
}
