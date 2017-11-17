package routes

import (
	"net/http"

	"bitbucket.org/exonch/ch-mail-templater/storages"
	"github.com/husobee/vestigo"
	"github.com/opentracing/opentracing-go"
)

type templateCreateRequest struct {
	Name    string `json:"template_name"`
	Version string `json:"template_version"`
	Data    string `json:"template_data"`
}

type templateCreateResponse struct {
	Name    string `json:"template_name"`
	Version string `json:"template_version"`
}

type templateUpdateRequest struct {
	Data string `json:"template_data"`
}

type templateUpdateResponse struct {
	Name    string `json:"template_name"`
	Version string `json:"template_version"`
}

type templateDeleteResponse struct {
	Name    string `json:"template_name"`
	Version string `json:"template_version"`
}

type templatesDeleteResponse struct {
	Name string `json:"template_name"`
}

func SetupTemplatesHandlers(router *vestigo.Router, tracer *opentracing.Tracer, storage *storages.TemplateStorage) {
	router.Post("/templates", templateCreateHandler,
		newOpenTracingMiddleware(tracer, "create template"),
		newTemplateStorageInjectionMiddleware(storage),
		newBodyUnmarshalMiddleware(templateCreateRequest{}))
	router.Get("/templates", templateGetHandler,
		newOpenTracingMiddleware(tracer, "retrieve template"),
		newTemplateStorageInjectionMiddleware(storage))
	router.Put("/templates/:template_name", templateUpdateHandler,
		newOpenTracingMiddleware(tracer, "update template"),
		newTemplateStorageInjectionMiddleware(storage),
		newBodyUnmarshalMiddleware(templateUpdateRequest{}))
	router.Delete("/templates/:template_name", templateDeleteHandler,
		newOpenTracingMiddleware(tracer, "delete template"),
		newTemplateStorageInjectionMiddleware(storage))
	// this is CRUD!!!
}

func templateCreateHandler(w http.ResponseWriter, r *http.Request) {
	storage := templateStorageFromContext(r.Context())
	request := bodyFromContext(r.Context()).(*templateCreateRequest)
	err := storage.PutTemplate(request.Name, request.Version, request.Data)
	if err != nil {
		log.WithError(err).Error("Create template failed")
		sendStorageError(w, err)
		return
	}
	sendJsonWithCode(w, http.StatusCreated, &templateCreateResponse{
		Name:    request.Name,
		Version: request.Version,
	})
}

func templateUpdateHandler(w http.ResponseWriter, r *http.Request) {
	storage := templateStorageFromContext(r.Context())
	data := bodyFromContext(r.Context()).(*templateUpdateRequest).Data
	name := vestigo.Param(r, "template_name")
	version := r.URL.Query().Get("version")
	err := storage.PutTemplate(name, version, data)
	if err != nil {
		log.WithError(err).Error("Update template failed")
		sendStorageError(w, err)
		return
	}
	sendJsonWithCode(w, http.StatusAccepted, &templateUpdateResponse{
		Name:    name,
		Version: version,
	})
}

func templateGetHandler(w http.ResponseWriter, r *http.Request) {
	storage := templateStorageFromContext(r.Context())
	name := vestigo.Param(r, "name")
	version := r.URL.Query().Get("version")
	var err error
	var respObj interface{}
	if version == "" { // if no "version" parameter specified, send all versions
		respObj, err = storage.GetTemplates(name)
	} else {
		respObj, err = storage.GetTemplate(name, version)
	}
	if err != nil {
		log.WithError(err).Error("Get template failed")
		sendStorageError(w, err)
		return
	}
	sendJson(w, respObj)
}

func templateDeleteHandler(w http.ResponseWriter, r *http.Request) {
	storage := templateStorageFromContext(r.Context())
	name := vestigo.Param(r, "name")
	version := r.URL.Query().Get("version")
	var err error
	var respObj interface{}
	if version == "" { // if no "version" parameter specified, delete all versions
		err = storage.DeleteTemplates(name)
		respObj = &templatesDeleteResponse{
			Name: name,
		}
	} else {
		err = storage.DeleteTemplate(name, version)
		respObj = &templateDeleteResponse{
			Name:    name,
			Version: version,
		}
	}
	if err != nil {
		log.WithError(err).Error("Delete template failed")
		sendStorageError(w, err)
		return
	}
	sendJson(w, respObj)
}
