package handlers

import (
	"net/http"

	"fmt"

	"git.containerum.net/ch/cherry"
	"git.containerum.net/ch/cherry/adaptors/gonic"
	"git.containerum.net/ch/mail-templater/pkg/models"
	"git.containerum.net/ch/mail-templater/pkg/mtErrors"
	m "git.containerum.net/ch/mail-templater/pkg/router/middleware"
	"git.containerum.net/ch/mail-templater/pkg/validation"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

//TemplateListGetHandler returns templates list
func TemplateListGetHandler(ctx *gin.Context) {
	svc := ctx.MustGet(m.MTServices).(*m.Services)

	respObj, err := svc.TemplateStorage.GetTemplatesList()
	if err != nil {
		if cherr, ok := err.(*cherry.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			ctx.Error(err)
			gonic.Gonic(mtErrors.ErrUnableGetMessagesList(), ctx)
		}
		return
	}
	ctx.JSON(http.StatusOK, respObj)
}

//TemplateGetHandler returns one template
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
		if cherr, ok := err.(*cherry.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			ctx.Error(err)
			gonic.Gonic(mtErrors.ErrUnableGetTemplate(), ctx)
		}
		return
	}
	ctx.JSON(http.StatusOK, respObj)
}

//TemplateCreateHandler creates template
func TemplateCreateHandler(ctx *gin.Context) {
	var request models.Template
	if err := ctx.ShouldBindWith(&request, binding.JSON); err != nil {
		gonic.Gonic(mtErrors.ErrRequestValidationFailed().AddDetailsErr(err), ctx)
		return
	}

	errs := validation.ValidateCreateTemplate(request)
	if errs != nil {
		gonic.Gonic(mtErrors.ErrRequestValidationFailed().AddDetailsErr(errs...), ctx)
		return
	}

	svc := ctx.MustGet(m.MTServices).(*m.Services)

	err := svc.TemplateStorage.PutTemplate(request.Name, request.Version, request.Data, request.Subject, true)
	if err != nil {
		if cherr, ok := err.(*cherry.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			ctx.Error(err)
			gonic.Gonic(mtErrors.ErrUnableSaveTemplate(), ctx)
		}
		return
	}
	ctx.JSON(http.StatusCreated, &models.Template{
		Name:    request.Name,
		Version: request.Version,
	})
}

//TemplateUpdateHandler updates template
func TemplateUpdateHandler(ctx *gin.Context) {
	var request models.Template
	if err := ctx.ShouldBindWith(&request, binding.JSON); err != nil {
		gonic.Gonic(mtErrors.ErrRequestValidationFailed().AddDetailsErr(err), ctx)
		return
	}

	errs := validation.ValidateUpdateTemplate(request)
	if errs != nil {
		gonic.Gonic(mtErrors.ErrRequestValidationFailed().AddDetailsErr(errs...), ctx)
		return
	}

	name := ctx.Param("name")
	version := ctx.Query("version")
	if version == "" {
		gonic.Gonic(mtErrors.ErrRequestValidationFailed().AddDetailsErr(fmt.Errorf(isRequiredQuery, "Version")), ctx)
		return
	}

	svc := ctx.MustGet(m.MTServices).(*m.Services)

	oldTemplate, err := svc.TemplateStorage.GetTemplate(name, version)
	if err != nil {
		if cherr, ok := err.(*cherry.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			ctx.Error(err)
			gonic.Gonic(mtErrors.ErrUnableUpdateTemplate(), ctx)
		}
		return
	}

	data := oldTemplate.Data
	subject := oldTemplate.Subject
	if request.Data != "" {
		data = request.Data
	}

	if request.Subject != "" {
		subject = request.Subject
	}

	err = svc.TemplateStorage.PutTemplate(name, version, data, subject, false)
	if err != nil {
		if cherr, ok := err.(*cherry.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			ctx.Error(err)
			gonic.Gonic(mtErrors.ErrUnableUpdateTemplate(), ctx)
		}
		return
	}
	ctx.JSON(http.StatusAccepted, &models.Template{
		Name:    name,
		Version: version,
	})
}

//TemplateDeleteHandler deletes template
func TemplateDeleteHandler(ctx *gin.Context) {
	name := ctx.Param("name")
	version, hasVersion := ctx.GetQuery("version")

	svc := ctx.MustGet(m.MTServices).(*m.Services)

	var err error
	var respObj interface{}
	if !hasVersion { // if no "version" parameter specified, delete all versions
		err = svc.TemplateStorage.DeleteTemplates(name)
		respObj = &models.Template{
			Name: name,
		}
	} else {
		err = svc.TemplateStorage.DeleteTemplate(name, version)
		respObj = &models.Template{
			Name:    name,
			Version: version,
		}
	}
	if err != nil {
		if cherr, ok := err.(*cherry.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			ctx.Error(err)
			gonic.Gonic(mtErrors.ErrUnableDeleteTemplate(), ctx)
		}
		return
	}
	ctx.JSON(http.StatusOK, respObj)
}
