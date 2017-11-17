package routes

import (
	"encoding/json"
	"net/http"

	"bitbucket.org/exonch/ch-mail-templater/storages"
	"golang.org/x/net/context"
)

func bodyFromContext(ctx context.Context) interface{} {
	return ctx.Value(bodyObjectContextKey)
}

func templateStorageFromContext(ctx context.Context) *storages.TemplateStorage {
	return ctx.Value(templateStorageContextKey).(*storages.TemplateStorage)
}

func messagesStorageFromContext(ctx context.Context) *storages.MessagesStorage {
	return ctx.Value(messagesStorageContextKey).(*storages.MessagesStorage)
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
	case storages.ErrTemplateNotExists, storages.ErrVersionNotExists, storages.ErrMessageNotExists:
		w.WriteHeader(http.StatusNotFound)
		if _, writeErr := w.Write([]byte(err.Error())); writeErr != nil {
			log.WithError(writeErr).Error("HTTP Response write error")
		}
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}
