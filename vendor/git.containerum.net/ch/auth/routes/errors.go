package routes

import (
	"encoding/json"
	"net/http"

	"git.containerum.net/ch/auth/storages"
	"github.com/sirupsen/logrus"
)

type httpResponseErrorBody struct {
	Error []string `json:"error"`
}

func sendErrorsWithCode(w http.ResponseWriter, errs []string, code int) {
	body, err := json.Marshal(&httpResponseErrorBody{
		Error: errs,
	})
	logrus.WithField("errors", errs).WithField("code", code).Debugf("Sending errors")
	if err != nil {
		logrus.Errorf("JSON Marshal: %v", err)
	}
	_, err = w.Write(body)
	if err != nil {
		logrus.Errorf("Response write: %v", err)
	}
	w.WriteHeader(code)
}

func sendError(w http.ResponseWriter, err error) {
	body, err := json.Marshal(&httpResponseErrorBody{
		Error: []string{err.Error()},
	})
	var code int

	switch err {
	case storages.ErrInvalidToken, storages.ErrTokenNotOwnedBySender:
		code = http.StatusUnauthorized
	default:
		code = http.StatusInternalServerError
	}

	logrus.WithField("error", err).WithField("code", code).Debugf("Sending error")
	_, err = w.Write(body)
	if err != nil {
		logrus.Errorf("Response write: %v", err)
	}
	w.WriteHeader(code)
}
