package handlers

import (
	"net/http"

	"fmt"

	mttypes "git.containerum.net/ch/json-types/mail-templater"
	ch "git.containerum.net/ch/kube-client/pkg/cherry"
	"git.containerum.net/ch/kube-client/pkg/cherry/adaptors/gonic"
	cherry "git.containerum.net/ch/kube-client/pkg/cherry/mail-templater"
	"git.containerum.net/ch/mail-templater/pkg/model"
	m "git.containerum.net/ch/mail-templater/pkg/router/middleware"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func TemplateCreateHandler(ctx *gin.Context) {
	var request mttypes.Template
	if err := ctx.ShouldBindWith(&request, binding.JSON); err != nil {
		gonic.Gonic(cherry.ErrRequestValidationFailed().AddDetailsErr(err), ctx)
		return
	}

	errs := model.ValidateCreateTemplate(request)
	if errs != nil {
		gonic.Gonic(cherry.ErrRequestValidationFailed().AddDetailsErr(errs...), ctx)
		return
	}

	svc := ctx.MustGet(m.MTServices).(*m.Services)

	err := svc.TemplateStorage.PutTemplate(request.Name, request.Version, request.Data, request.Subject, true)
	if err != nil {
		if cherr, ok := err.(*ch.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			ctx.Error(err)
			gonic.Gonic(cherry.ErrUnableSaveTemplate(), ctx)
		}
		return
	}
	ctx.JSON(http.StatusCreated, &mttypes.Template{
		Name:    request.Name,
		Version: request.Version,
	})
}

func TemplateUpdateHandler(ctx *gin.Context) {
	var request mttypes.Template
	if err := ctx.ShouldBindWith(&request, binding.JSON); err != nil {
		gonic.Gonic(cherry.ErrRequestValidationFailed().AddDetailsErr(err), ctx)
		return
	}

	errs := model.ValidateUpdateTemplate(request)
	if errs != nil {
		gonic.Gonic(cherry.ErrRequestValidationFailed().AddDetailsErr(errs...), ctx)
		return
	}

	name := ctx.Param("name")
	version := ctx.Query("version")
	if version == "" {
		gonic.Gonic(cherry.ErrRequestValidationFailed().AddDetailsErr(fmt.Errorf(model.IsRequiredQuery, "Version")), ctx)
		return
	}

	svc := ctx.MustGet(m.MTServices).(*m.Services)

	oldTemplate, err := svc.TemplateStorage.GetTemplate(name, version)
	if err != nil {
		if cherr, ok := err.(*ch.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			ctx.Error(err)
			gonic.Gonic(cherry.ErrUnableUpdateTemplate(), ctx)
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
		if cherr, ok := err.(*ch.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			if cherr, ok := err.(*ch.Err); ok {
				gonic.Gonic(cherr, ctx)
			} else {
				ctx.Error(err)
				gonic.Gonic(cherry.ErrUnableUpdateTemplate(), ctx)
			}
		}
		return
	}
	ctx.JSON(http.StatusAccepted, &mttypes.Template{
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
		if cherr, ok := err.(*ch.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			ctx.Error(err)
			gonic.Gonic(cherry.ErrUnableGetTemplate(), ctx)
		}
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
		respObj = &mttypes.Template{
			Name: name,
		}
	} else {
		err = svc.TemplateStorage.DeleteTemplate(name, version)
		respObj = &mttypes.Template{
			Name:    name,
			Version: version,
		}
	}
	if err != nil {
		if cherr, ok := err.(*ch.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			ctx.Error(err)
			gonic.Gonic(cherry.ErrUnableDeleteTemplate(), ctx)
		}
		return
	}
	ctx.JSON(http.StatusOK, respObj)
}

func TemplateListGetHandler(ctx *gin.Context) {
	svc := ctx.MustGet(m.MTServices).(*m.Services)

	respObj, err := svc.TemplateStorage.GetTemplatesList()
	if err != nil {
		if cherr, ok := err.(*ch.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			ctx.Error(err)
			gonic.Gonic(cherry.ErrUnableGetMessagesList(), ctx)
		}
		return
	}
	ctx.JSON(http.StatusOK, respObj)
}
