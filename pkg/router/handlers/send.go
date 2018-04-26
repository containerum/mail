package handlers

import (
	"net/http"

	"git.containerum.net/ch/cherry"
	"git.containerum.net/ch/cherry/adaptors/gonic"
	"git.containerum.net/ch/mail-templater/pkg/models"
	"git.containerum.net/ch/mail-templater/pkg/mtErrors"
	m "git.containerum.net/ch/mail-templater/pkg/router/middleware"
	"git.containerum.net/ch/mail-templater/pkg/validation"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

//SimpleSendHandler sends email in simple way
func SimpleSendHandler(ctx *gin.Context) {
	var request models.SimpleSendRequest
	if err := ctx.ShouldBindWith(&request, binding.JSON); err != nil {
		gonic.Gonic(mtErrors.ErrRequestValidationFailed().AddDetailsErr(err), ctx)
		return
	}

	errs := validation.ValidateSimpleSendRequest(request)
	if errs != nil {
		gonic.Gonic(mtErrors.ErrRequestValidationFailed().AddDetailsErr(errs...), ctx)
		return
	}

	svc := ctx.MustGet(m.MTServices).(*m.Services)

	_, tv, err := svc.TemplateStorage.GetLatestVersionTemplate(request.Template)
	if err != nil {
		ctx.Error(err)
		gonic.Gonic(mtErrors.ErrTemplateNotExist(), ctx)
		return
	}

	info, err := svc.UserManagerClient.UserInfoByID(ctx, request.UserID)
	if err != nil {
		ctx.Error(err)
		cherr, ok := err.(*cherry.Err)
		if ok {
			gonic.Gonic(cherr, ctx)
		} else {
			gonic.Gonic(mtErrors.ErrMailSendFailed().AddDetailsErr(err), ctx)
		}
		return
	}

	recipient := &models.Recipient{
		ID:        request.UserID,
		Name:      info.Login,
		Email:     info.Login,
		Variables: request.Variables,
	}

	status, err := svc.UpstreamSimple.SimpleSend(ctx, request.Template, tv, recipient)
	if err != nil {
		ctx.Error(err)
		gonic.Gonic(mtErrors.ErrMailSendFailed(), ctx)
		return
	}
	ctx.JSON(http.StatusOK, models.SimpleSendResponse{
		UserID: status.RecipientID,
	})
}

//SendHandler sends email in not so simple way
func SendHandler(ctx *gin.Context) {
	name := ctx.Param("name")
	version, hasVersion := ctx.GetQuery("version")
	var request models.SendRequest
	if err := ctx.ShouldBindWith(&request, binding.JSON); err != nil {
		gonic.Gonic(mtErrors.ErrRequestValidationFailed().AddDetailsErr(err), ctx)
		return
	}

	errs := validation.ValidateSendRequest(request)
	if errs != nil {
		gonic.Gonic(mtErrors.ErrRequestValidationFailed().AddDetailsErr(errs...), ctx)
		return
	}

	svc := ctx.MustGet(m.MTServices).(*m.Services)

	var tv *models.Template
	var err error
	if !hasVersion {
		_, tv, err = svc.TemplateStorage.GetLatestVersionTemplate(name)
	} else {
		tv, err = svc.TemplateStorage.GetTemplate(name, version)
	}
	if err != nil {
		if cherr, ok := err.(*cherry.Err); ok {
			gonic.Gonic(cherr, ctx)
		} else {
			ctx.Error(err)
			gonic.Gonic(mtErrors.ErrMailSendFailed(), ctx)
		}
		return
	}
	status, err := svc.Upstream.Send(ctx, name, tv, &request)
	if err != nil {
		ctx.Error(err)
		gonic.Gonic(mtErrors.ErrMailSendFailed().AddDetailsErr(err), ctx)
		return
	}
	ctx.JSON(http.StatusOK, status)
}
