package routes

import (
	"net/http"

	"bitbucket.org/exonch/ch-mail-templater/storages"
	"github.com/husobee/vestigo"
	"github.com/opentracing/opentracing-go"
)

type messageGetResponse struct {
	Id string `json:"id"`
	*storages.MessagesStorageValue
}

func SetupMessagesHandlers(router *vestigo.Router, tracer *opentracing.Tracer,
	storage *storages.MessagesStorage /*TODO: Upstream*/) {
	router.Get("/messages/:message_id", messageGetHandler,
		newOpenTracingMiddleware(tracer, "get message copy"),
		newMessagesStorageInjectionMiddleware(storage))
}

func messageGetHandler(w http.ResponseWriter, r *http.Request) {
	storage := messagesStorageFromContext(r.Context())
	id := vestigo.Param(r, "message_id")
	v, err := storage.GetValue(id)
	if err != nil {
		log.WithError(err).Error("Get message failed")
		sendStorageError(w, err)
		return
	}
	sendJson(w, &messageGetResponse{
		Id:                   id,
		MessagesStorageValue: v,
	})
}
