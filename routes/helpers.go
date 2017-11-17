package routes

import (
	"encoding/json"
	"net/http"

	templater "bitbucket.org/exonch/ch-mail-templater"
	"golang.org/x/net/context"
)

func bodyFromContext(ctx context.Context) interface{} {
	return ctx.Value(bodyObjectContextKey)
}

func storageFromContext(ctx context.Context) *templater.TemplateStorage {
	return ctx.Value(storageContextKey).(*templater.TemplateStorage)
}

func sendJsonWithCode(w http.ResponseWriter, code int, resp interface{}) {
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.WithError(err).Error("Response write failed")
	}
}

func sendJson(w http.ResponseWriter, resp interface{}) {
	sendJsonWithCode(w, http.StatusOK, resp)
}

func sendStorageError(w http.ResponseWriter, err error) {
	log.WithError(err).Debugf("Sending storage error")
	switch err {
	case nil:
	case templater.ErrTemplateNotExists, templater.ErrVersionNotExists:
		w.WriteHeader(http.StatusNotFound)
		if _, writeErr := w.Write([]byte(err.Error())); writeErr != nil {
			log.WithError(writeErr).Error("HTTP Response write error")
		}
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}
