package routes

import (
	"net/http"

	"io/ioutil"

	"encoding/json"

	templater "bitbucket.org/exonch/ch-mail-templater"
	"github.com/husobee/vestigo"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

const (
	storageContextKey    = "storage"
	upstreamContextKey   = "upstream"
	bodyObjectContextKey = "body"
)

var log = logrus.WithField("component", "http_router")

// Middleware for opentracing functionality. MUST BE FIRST in chain
func newOpenTracingMiddleware(tracer opentracing.Tracer, operationName string) vestigo.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			log.WithField("request", r).Debugf("Opentracing middleware")
			wireContext, err := tracer.Extract(
				opentracing.TextMap,
				opentracing.HTTPHeadersCarrier(r.Header),
			)
			if err != nil {
				log.Errorf("Opentracing span extract: %v", err)
			}

			span := tracer.StartSpan(operationName, ext.RPCServerOption(wireContext))
			defer span.Finish()

			ctx := opentracing.ContextWithSpan(r.Context(), span)
			next(w, r.WithContext(ctx))
		}
	}
}

// Middleware injecting storage interface to context. MUST BE INCLUDED if storage used in handler
func newStorageInjectionMiddleware(storage *templater.TemplateStorage) vestigo.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			log.WithField("request", r).Debug("StorageInjection middleware")
			ctx := context.WithValue(r.Context(), storageContextKey, storage)
			next(w, r.WithContext(ctx))
		}
	}
}

// func newUpstreamInjectionMiddleware()

func newBodyUnmarshalMiddleware(obj interface{}) vestigo.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			logger := log.WithField("request", r)

			body, err := ioutil.ReadAll(r.Body)
			logger = logger.WithField("body", string(body))

			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
			}
			r.Body.Close()

			if err := json.Unmarshal(body, &obj); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			r.WithContext(context.WithValue(r.Context(), bodyObjectContextKey, &obj))
			next(w, r)
		}
	}
}
