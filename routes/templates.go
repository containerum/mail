package routes

import (
	"net/http"

	"git.containerum.net/ch/mail-templater/upstreams"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type templateCreateRequest struct {
	Name    string `json:"template_name" binding:"required"`
	Version string `json:"template_version" binding:"required"`
	Data    string `json:"template_data" binding:"required,base64"`
	Subject string `json:"template_subject" binding:"required"`
}

type templateCreateResponse struct {
	Name    string `json:"template_name"`
	Version string `json:"template_version"`
}

type templateUpdateRequest struct {
	Data    string `json:"template_data" binding:"required,base64"`
	Subject string `json:"template_subject" binding:"required"`
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

func templateCreateHandler(ctx *gin.Context) {
	var request templateCreateRequest
	if err := ctx.ShouldBindWith(&request, binding.JSON); err != nil {
		ctx.Error(err)
		sendValidationError(ctx, err)
		return
	}
	err := svc.TemplateStorage.PutTemplate(request.Name, request.Version, request.Data, request.Subject)
	if err != nil {
		ctx.Error(err)
		sendStorageError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, &templateCreateResponse{
		Name:    request.Name,
		Version: request.Version,
	})
}

func templateUpdateHandler(ctx *gin.Context) {
	var request templateUpdateRequest
	if err := ctx.ShouldBindWith(&request, binding.JSON); err != nil {
		ctx.Error(err)
		sendValidationError(ctx, err)
		return
	}
	name := ctx.Param("template_name")
	version := ctx.Query("version")
	err := svc.TemplateStorage.PutTemplate(name, version, request.Data, request.Subject)
	if err != nil {
		ctx.Error(err)
		sendStorageError(ctx, err)
		return
	}
	ctx.JSON(http.StatusAccepted, &templateUpdateResponse{
		Name:    name,
		Version: version,
	})
}

func templateGetHandler(ctx *gin.Context) {
	name := ctx.Param("name")
	version, hasVersion := ctx.GetQuery("version")
	var err error
	var respObj interface{}
	if !hasVersion { // if no "version" parameter specified, send all versions
		respObj, err = svc.TemplateStorage.GetTemplates(name)
	} else {
		respObj, err = svc.TemplateStorage.GetTemplate(name, version)
	}
	if err != nil {
		ctx.Error(err)
		sendStorageError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, respObj)
}

func templateDeleteHandler(ctx *gin.Context) {
	name := ctx.Param("name")
	version, hasVersion := ctx.GetQuery("version")
	var err error
	var respObj interface{}
	if !hasVersion { // if no "version" parameter specified, delete all versions
		err = svc.TemplateStorage.DeleteTemplates(name)
		respObj = &templatesDeleteResponse{
			Name: name,
		}
	} else {
		err = svc.TemplateStorage.DeleteTemplate(name, version)
		respObj = &templateDeleteResponse{
			Name:    name,
			Version: version,
		}
	}
	if err != nil {
		ctx.Error(err)
		sendStorageError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, respObj)
}

func templateSendHandler(ctx *gin.Context) {
	name := ctx.Param("name")
	version := ctx.Query("version")
	var request upstreams.SendRequest
	if err := ctx.ShouldBindWith(&request, binding.JSON); err != nil {
		ctx.Error(err)
		sendValidationError(ctx, err)
		return
	}
	tv, err := svc.TemplateStorage.GetTemplate(name, version)
	if err != nil {
		ctx.Error(err)
		sendStorageError(ctx, err)
		return
	}
	status, err := svc.Upstream.Send(name, tv, &request)
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, status)
}
