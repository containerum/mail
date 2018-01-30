package routes

import (
	"net/http"

	mttypes "git.containerum.net/ch/json-types/mail-templater"
	"github.com/gin-gonic/gin"
)

func templateCreateHandler(ctx *gin.Context) {
	var request mttypes.TemplateCreateRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
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
	ctx.JSON(http.StatusCreated, &mttypes.TemplateCreateResponse{
		Name:    request.Name,
		Version: request.Version,
	})
}

func templateUpdateHandler(ctx *gin.Context) {
	var request mttypes.TemplateUpdateRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.Error(err)
		sendValidationError(ctx, err)
		return
	}
	name := ctx.Param("name")
	version := ctx.Query("version")

	respObj, err := svc.TemplateStorage.GetTemplate(name, version)
	if err != nil {
		ctx.Error(err)
		sendStorageError(ctx, err)
		return
	}

	data := respObj.Data
	subject := respObj.Subject
	if request.Data != "" {
		data = request.Data
	}

	if request.Subject != "" {
		subject = request.Subject
	}

	err = svc.TemplateStorage.PutTemplate(name, version, data, subject)
	if err != nil {
		ctx.Error(err)
		sendStorageError(ctx, err)
		return
	}
	ctx.JSON(http.StatusAccepted, &mttypes.TemplateUpdateResponse{
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
		respObj = &mttypes.TemplatesDeleteResponse{
			Name: name,
		}
	} else {
		err = svc.TemplateStorage.DeleteTemplate(name, version)
		respObj = &mttypes.TemplateDeleteResponse{
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
	version, hasVersion := ctx.GetQuery("version")
	var request mttypes.SendRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.Error(err)
		sendValidationError(ctx, err)
		return
	}
	var tv *mttypes.TemplateStorageValue
	var err error
	if !hasVersion {
		_, tv, err = svc.TemplateStorage.GetLatestVersionTemplate(name)
	} else {
		tv, err = svc.TemplateStorage.GetTemplate(name, version)
	}
	if err != nil {
		ctx.Error(err)
		sendStorageError(ctx, err)
		return
	}
	status, err := svc.Upstream.Send(ctx, name, tv, &request)
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, status)
}

func templateListGetHandler(ctx *gin.Context) {
	respObj, err := svc.TemplateStorage.GetTemplatesList()
	if err != nil {
		ctx.Error(err)
		sendStorageError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, respObj)
}
