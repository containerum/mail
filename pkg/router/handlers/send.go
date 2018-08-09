package handlers

import (
	"net/http"

	"git.containerum.net/ch/mail-templater/pkg/models"
	"git.containerum.net/ch/mail-templater/pkg/mterrors"
	m "git.containerum.net/ch/mail-templater/pkg/router/middleware"
	"git.containerum.net/ch/mail-templater/pkg/validation"
	"github.com/containerum/cherry"
	"github.com/containerum/cherry/adaptors/gonic"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// swagger:operation POST /send Send SimpleSendHandler
// Send message from internal service.
// https://ch.pages.containerum.net/api-docs/modules/ch-mail-template/index.html#send-message-for-resource-manager
//
// ---
// x-method-visibility: public
// parameters:
//  - name: body
//    in: body
//    schema:
//      $ref: '#/definitions/SimpleSendRequest'
// responses:
//  '202':
//    description: message simple send response
//    schema:
//      $ref: '#/definitions/SimpleSendResponse'
//  default:
//    $ref: '#/responses/error'
func SimpleSendHandler(ctx *gin.Context) {
	var request models.SimpleSendRequest
	if err := ctx.ShouldBindWith(&request, binding.JSON); err != nil {
		gonic.Gonic(mterrors.ErrRequestValidationFailed().AddDetailsErr(err), ctx)
		return
	}

	errs := validation.ValidateSimpleSendRequest(request)
	if errs != nil {
		gonic.Gonic(mterrors.ErrRequestValidationFailed().AddDetailsErr(errs...), ctx)
		return
	}

	svc := ctx.MustGet(m.MTServices).(*m.Services)

	_, tv, err := svc.TemplateStorage.GetLatestVersionTemplate(request.Template)
	if err != nil {
		ctx.Error(err)
		gonic.Gonic(mterrors.ErrTemplateNotExist(), ctx)
		return
	}

	info, err := svc.UserManagerClient.UserInfoByID(ctx, request.UserID)
	if err != nil {
		ctx.Error(err)
		cherr, ok := err.(*cherry.Err)
		if ok {
			gonic.Gonic(cherr, ctx)
		} else {
			gonic.Gonic(mterrors.ErrMailSendFailed().AddDetailsErr(err), ctx)
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
		gonic.Gonic(mterrors.ErrMailSendFailed(), ctx)
		return
	}
	ctx.JSON(http.StatusAccepted, models.SimpleSendResponse{
		UserID: status.RecipientID,
	})
}

// swagger:operation POST /templates/{template} Send SendHandler
// Send message to any email.
// https://ch.pages.containerum.net/api-docs/modules/ch-mail-template/index.html#send-message-extended
//
// ---
// x-method-visibility: public
// parameters:
//  - name: template
//    in: path
//    type: string
//    required: true
// parameters:
//  - name: template
//    in: path
//    type: string
//    required: true
//  - name: body
//    in: body
//    schema:
//      $ref: '#/definitions/SendRequest'
// responses:
//  '202':
//    description: message send response
//    schema:
//      $ref: '#/definitions/SendResponse'
//  default:
//    $ref: '#/responses/error'
func SendHandler(ctx *gin.Context) {
	name := ctx.Param("name")
	version, hasVersion := ctx.GetQuery("version")
	var request models.SendRequest
	if err := ctx.ShouldBindWith(&request, binding.JSON); err != nil {
		gonic.Gonic(mterrors.ErrRequestValidationFailed().AddDetailsErr(err), ctx)
		return
	}

	errs := validation.ValidateSendRequest(request)
	if errs != nil {
		gonic.Gonic(mterrors.ErrRequestValidationFailed().AddDetailsErr(errs...), ctx)
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
			gonic.Gonic(mterrors.ErrMailSendFailed(), ctx)
		}
		return
	}
	status, err := svc.Upstream.Send(ctx, name, tv, &request)
	if err != nil {
		ctx.Error(err)
		gonic.Gonic(mterrors.ErrMailSendFailed().AddDetailsErr(err), ctx)
		return
	}
	ctx.JSON(http.StatusAccepted, status)
}
