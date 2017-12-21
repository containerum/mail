package routes

import (
	"net/http"

	"bytes"
	"fmt"
	"io/ioutil"

	"git.containerum.net/ch/grpc-proto-files/auth"
	"github.com/husobee/vestigo"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// Middleware for opentracing functionality. MUST BE FIRST in chain
func newOpenTracingMiddleware(tracer opentracing.Tracer, operationName string) vestigo.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			logrus.WithField("request", r).Debugf("Opentracing middleware")
			wireContext, err := tracer.Extract(
				opentracing.TextMap,
				opentracing.HTTPHeadersCarrier(r.Header),
			)
			if err != nil {
				logrus.Errorf("Opentracing span extract: %v", err)
			}

			span := tracer.StartSpan(operationName, ext.RPCServerOption(wireContext))
			defer span.Finish()

			ctx := opentracing.ContextWithSpan(r.Context(), span)
			next(w, r.WithContext(ctx))
		}
	}
}

// Middleware injecting storage interface to context. MUST BE INCLUDED if storage used in handler
func newStorageInjectionMiddleware(storage auth.AuthServer) vestigo.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			logrus.WithField("request", r).Debug("StorageInjection middleware")
			ctx := context.WithValue(r.Context(), authServerContextKey, storage)
			next(w, r.WithContext(ctx))
		}
	}
}

// name -> function validating value
type validators map[string](func(value string) error)

func newHeaderValidationMiddleware(validators validators) vestigo.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			logger := logrus.WithField("request", r)

			logger.Debugf("HeaderValidation middleware")
			validationErrors := make(map[string]error) // header name to error

			for headerName, validator := range validators {
				headerValue := r.Header.Get(headerName)
				if headerValue != "" {
					logger.Debugf("Validating header %s: %s", headerName, headerValue)
					if err := validator(headerValue); err != nil {
						validationErrors[headerName] = err
					}
				}
			}

			if len(validationErrors) != 0 {
				var errs []string
				for header, err := range validationErrors {
					errs = append(errs, fmt.Sprintf("Invalid header %s: %v", header, err))
				}
				sendErrorsWithCode(w, errs, http.StatusBadRequest)
				return
			}

			next(w, r)
		}
	}
}

func newParameterValidationMiddleware(validators validators) vestigo.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			logger := logrus.WithField("request", r)

			logger.Debugf("ParameterValidation middleware")
			validationErrors := make(map[string]error) // param name to error

			for paramName, validator := range validators {
				paramValue := vestigo.Param(r, paramName)
				if paramValue != "" {
					logger.Debugf("Validating parameter %s: %s", paramName, paramValue)
					if err := validator(paramValue); err != nil {
						validationErrors[paramName] = err
					}
				}
			}

			if len(validationErrors) != 0 {
				var errs []string
				for header, err := range validationErrors {
					errs = append(errs, fmt.Sprintf("Invalid parameter %s: %v\n", header, err))
				}
				sendErrorsWithCode(w, errs, http.StatusBadRequest)
				return
			}

			next(w, r)
		}
	}
}

func newBodyValidationMiddleware(validator func(body []byte) error) vestigo.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			logger := logrus.WithField("request", r)

			body, err := ioutil.ReadAll(r.Body)
			logger = logger.WithField("body", string(body))

			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
			}
			r.Body.Close()

			if err := validator(body); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			next(w, r)
		}
	}
}
