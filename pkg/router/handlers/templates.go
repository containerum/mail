package handlers

import (
	"net/http"

	"fmt"

	"git.containerum.net/ch/mail-templater/pkg/models"
	"git.containerum.net/ch/mail-templater/pkg/mtErrors"
	m "git.containerum.net/ch/mail-templater/pkg/router/middleware"
	"git.containerum.net/ch/mail-templater/pkg/validation"
	"github.com/containerum/cherry"
	"github.com/containerum/cherry/adaptors/gonic"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// swagger:operation GET /templates Templates TemplateListGetHandler
// Get templates list.
// https://ch.pages.containerum.net/api-docs/modules/ch-mail-template/index.html#get-all-templates
//
// ---
// x-method-visibility: private
// parameters:
//  - $ref: '#/parameters/UserRoleHeader'
// responses:
//  '200':
//    description: templates list get response
//    schema:
//      $ref: '#/definitions/TemplatesListResponse'
//  default:
//    $ref: '#/responses/error'
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

// swagger:operation GET /templates/{name} Templates TemplateGetHandler
// Get single template.
// https://ch.pages.containerum.net/api-docs/modules/ch-mail-template/index.html#get-all-versions-of-template
// https://ch.pages.containerum.net/api-docs/modules/ch-mail-template/index.html#get-template-of-specific-version
//
// ---
// x-method-visibility: private
// parameters:
//  - $ref: '#/parameters/UserRoleHeader'
//  - name: version
//    in: query
//    type: string
//    required: false
//  - name: name
//    in: path
//    type: string
//    required: true
// responses:
//  '200':
//    description: templates list get response
//    schema:
//      $ref: '#/definitions/TemplatesListResponse'
//  default:
//    $ref: '#/responses/error'
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

// swagger:operation POST /templates Templates TemplateCreateHandler
// Create new template.
// https://ch.pages.containerum.net/api-docs/modules/ch-mail-template/index.html#create-template
//
// ---
// x-method-visibility: private
// parameters:
//  - $ref: '#/parameters/UserRoleHeader'
//  - name: body
//    in: body
//    schema:
//      $ref: '#/definitions/Template'
// responses:
//  '201':
//    description: created template
//    schema:
//      $ref: '#/definitions/Template'
//  default:
//    $ref: '#/responses/error'
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

// swagger:operation PUT /templates/{name} Templates TemplateUpdateHandler
// Update template.
// https://ch.pages.containerum.net/api-docs/modules/ch-mail-template/index.html#update-template-of-specific-version
//
// ---
// x-method-visibility: private
// parameters:
//  - $ref: '#/parameters/UserRoleHeader'
//  - name: name
//    in: path
//    type: string
//    required: true
//  - name: body
//    in: body
//    schema:
//      $ref: '#/definitions/Template'
// responses:
//  '202':
//    description: updated template
//    schema:
//      $ref: '#/definitions/Template'
//  default:
//    $ref: '#/responses/error'
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

// swagger:operation DELETE /templates/{name} Templates TemplateDeleteHandler
// Delete template.
// https://ch.pages.containerum.net/api-docs/modules/ch-mail-template/index.html#update-template-of-specific-version
//
// ---
// x-method-visibility: private
// parameters:
//  - $ref: '#/parameters/UserRoleHeader'
//  - name: name
//    in: path
//    type: string
//    required: true
//  - name: version
//    in: query
//    type: string
//    required: false
// responses:
//  '202':
//    description: template deleted
//  default:
//    $ref: '#/responses/error'
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
	ctx.JSON(http.StatusAccepted, respObj)
}
