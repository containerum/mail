package handlers

import (
	"net/http"

	"git.containerum.net/ch/json-types/errors"
	mttypes "git.containerum.net/ch/json-types/mail-templater"
	m "git.containerum.net/ch/mail-templater/pkg/router/middleware"
	"github.com/gin-gonic/gin"
)

func TemplateCreateHandler(ctx *gin.Context) {
	var request mttypes.TemplateCreateRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.Error(err)
		//	ctx.AbortWithStatusJSON(http.StatusBadRequest, ParseBindErorrs(err))
		return
	}

	svc := ctx.MustGet(m.MTServices).(*m.Services)

	err := svc.TemplateStorage.PutTemplate(request.Name, request.Version, request.Data, request.Subject)
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatusJSON(errors.ErrorWithHTTPStatus(err))
		return
	}
	ctx.JSON(http.StatusCreated, &mttypes.TemplateCreateResponse{
		Name:    request.Name,
		Version: request.Version,
	})
}

func TemplateUpdateHandler(ctx *gin.Context) {
	var request mttypes.TemplateUpdateRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.Error(err)
		//	ctx.AbortWithStatusJSON(http.StatusBadRequest, ParseBindErorrs(err))
		return
	}
	name := ctx.Param("name")
	version := ctx.Query("version")

	svc := ctx.MustGet(m.MTServices).(*m.Services)

	respObj, err := svc.TemplateStorage.GetTemplate(name, version)
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatusJSON(errors.ErrorWithHTTPStatus(err))
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
		ctx.AbortWithStatusJSON(errors.ErrorWithHTTPStatus(err))
		return
	}
	ctx.JSON(http.StatusAccepted, &mttypes.TemplateUpdateResponse{
		Name:    name,
		Version: version,
	})
}

func TemplateGetHandler(ctx *gin.Context) {
	name := ctx.Param("name")
	version, hasVersion := ctx.GetQuery("version")

	svc := ctx.MustGet(m.MTServices).(*m.Services)

	var err error
	var respObj interface{}
	if !hasVersion { // if no "version" parameter specified, send all versions
		respObj, err = svc.TemplateStorage.GetTemplates(name)
	} else {
		respObj, err = svc.TemplateStorage.GetTemplate(name, version)
	}
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatusJSON(errors.ErrorWithHTTPStatus(err))
		return
	}
	ctx.JSON(http.StatusOK, respObj)
}

func TemplateDeleteHandler(ctx *gin.Context) {
	name := ctx.Param("name")
	version, hasVersion := ctx.GetQuery("version")

	svc := ctx.MustGet(m.MTServices).(*m.Services)

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
		ctx.AbortWithStatusJSON(errors.ErrorWithHTTPStatus(err))
		return
	}
	ctx.JSON(http.StatusOK, respObj)
}

func TemplateSendHandler(ctx *gin.Context) {
	name := ctx.Param("name")
	version, hasVersion := ctx.GetQuery("version")
	var request mttypes.SendRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.Error(err)
		//	ctx.AbortWithStatusJSON(http.StatusBadRequest, ParseBindErorrs(err))
		return
	}

	svc := ctx.MustGet(m.MTServices).(*m.Services)

	var tv *mttypes.TemplateStorageValue
	var err error
	if !hasVersion {
		_, tv, err = svc.TemplateStorage.GetLatestVersionTemplate(name)
	} else {
		tv, err = svc.TemplateStorage.GetTemplate(name, version)
	}
	if err != nil {
		ctx.Error(err)
		ctx.AbortWithStatusJSON(errors.ErrorWithHTTPStatus(err))
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

func TemplateListGetHandler(ctx *gin.Context) {
	svc := ctx.MustGet(m.MTServices).(*m.Services)

	respObj, err := svc.TemplateStorage.GetTemplatesList()
	if err != nil {
		ctx.Error(err)
		//sendStorageError(ctx, err)
		ctx.AbortWithStatusJSON(errors.ErrorWithHTTPStatus(err))
		return
	}
	ctx.JSON(http.StatusOK, respObj)
}
