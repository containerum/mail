package rserrors

import (
	"net/http"

	"git.containerum.net/ch/kube-client/pkg/cherry"
)

var resourceErr = cherry.BuildErr(cherry.ResourceService)

func ErrDatabase() *cherry.Err {
	return resourceErr("Database error", http.StatusInternalServerError, 1)
}
func ErrResourceNotExists() *cherry.Err {
	return resourceErr("Resource is not exists", http.StatusNotFound, 2)
}
func ErrResourceAlreadyExists() *cherry.Err {
	return resourceErr("Resource already exists", http.StatusConflict, 3)
}
func ErrPermissionDenied() *cherry.Err {
	return resourceErr("Permission denied", http.StatusForbidden, 4)
}
func ErrTariffUnchanged() *cherry.Err {
	return resourceErr("Tariff unchanged", http.StatusBadRequest, 5)
}
func ErrTariffNotFound() *cherry.Err {
	return resourceErr("Tariff was not found", http.StatusNotFound, 6)
}
func ErrResourceNotOwned() *cherry.Err {
	return resourceErr("Can`t set access for resource which not owned by user", http.StatusForbidden, 7)
}
func ErrDeleteOwnerAccess() *cherry.Err {
	return resourceErr("Owner can`t delete has own access to resource", http.StatusConflict, 8)
}
func ErrAccessRecordNotExists() *cherry.Err {
	return resourceErr("Access record for user not exists", http.StatusNotFound, 9)
}
func ErrInternal() *cherry.Err {
	return resourceErr("Internal error", http.StatusInternalServerError, 10)
}
func ErrValidation() *cherry.Err {
	return resourceErr("Validation error", http.StatusBadRequest, 11)
}
func ErrServiceNotExternal() *cherry.Err {
	return resourceErr("Service is not external", http.StatusBadRequest, 12)
}
func ErrTCPPortNotFound() *cherry.Err {
	return resourceErr("TCP Port was not found in service", http.StatusNotFound, 13)
}
