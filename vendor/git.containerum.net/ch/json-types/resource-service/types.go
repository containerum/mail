package resource

import (
	"time"

	"database/sql"

	"git.containerum.net/ch/json-types/misc"
)

type Kind string // constants KindNamespace, KindVolume, ... It`s recommended to use strings.ToLower before comparsion

const (
	KindNamespace  Kind = "namespace"
	KindVolume          = "volume"
	KindExtService      = "extservice"
	KindIntService      = "intservice"
	KindDomain          = "domain"
)

type PermissionStatus string // constants PermissionStatusOwner, PermissionStatusRead

const (
	PermissionStatusOwner      PermissionStatus = "owner"
	PermissionStatusRead                        = "read"
	PermissionStatusWrite                       = "write"
	PermissionStatusReadDelete                  = "readdelete"
	PermissionStatusNone                        = "none"
)

type Resource struct {
	ID         string          `json:"id" db:"id"`
	CreateTime string          `json:"create_time,omitempty" db:"create_time"`
	Deleted    bool            `json:"deleted,omitempty" db:"deleted"`
	DeleteTime misc.PqNullTime `json:"delete_time,omitempty" db:"delete_time"`
	TariffID   string          `json:"tariff_id,omitempty" db:"tariff_id"`
}

type Namespace struct {
	Resource

	RAM                 int `json:"ram" db:"ram"` // megabytes
	CPU                 int `json:"cpu" db:"cpu"`
	MaxExternalServices int `json:"max_external_services" db:"max_ext_services"`
	MaxIntServices      int `json:"max_internal_services" db:"max_int_services"`
	MaxTraffic          int `json:"max_traffic" db:"max_traffic"` // megabytes per month
}

type Volume struct {
	Resource

	Active     *bool `json:"active,omitempty" db:"active"`
	Capacity   int   `json:"capacity" db:"capacity"` // gigabytes
	Replicas   int   `json:"replicas" db:"replicas"`
	Persistent bool  `json:"is_persistent" db:"is_persistent"`
}

type Deployment struct {
	ID          string          `json:"id" db:"id"`
	NamespaceID string          `json:"namespace_id" db:"ns_id"`
	Name        string          `json:"name" db:"name"`
	RAM         int             `json:"ram" db:"ram"`
	CPU         int             `json:"cpu" db:"cpu"`
	CreateTime  time.Time       `json:"create_time,omitempty" db:"create_time"`
	Deleted     *bool           `json:"deleted,omitempty" db:"deleted"`
	DeleteTime  misc.PqNullTime `json:"delete_time,omitempty" db:"delete_time"`
}

type PermissionRecord struct {
	ID                    string           `json:"id,omitempty" db:"id"`
	Kind                  Kind             `json:"kind,omitempty" db:"kind"`
	ResourceID            sql.NullString   `json:"resource_id,omitempty" db:"resource_id"` // it can be null for resources without tables
	ResourceLabel         string           `json:"label,omitempty" db:"resource_label"`
	OwnerUserID           string           `json:"owner_user_id,omitempty" db:"owner_user_id"`
	CreateTime            time.Time        `json:"create_time,omitempty" db:"create_time"`
	UserID                string           `json:"user_id" db:"user_id"`
	AccessLevel           PermissionStatus `json:"access_level" db:"access_level"`
	Limited               bool             `json:"limited" db:"limited"`
	AccessLevelChangeTime time.Time        `json:"access_level_change_time" db:"access_level_change_time"`
	NewAccessLevel        PermissionStatus `json:"new_access_level,omitempty" db:"new_access_level"`
}

// Types below is not for storing in db

type NamespaceWithPermission struct {
	Namespace
	PermissionRecord
}

type VolumeWithPermission struct {
	Volume
	PermissionRecord
}

type NamespaceWithVolumes struct {
	NamespaceWithPermission
	Volume []VolumeWithPermission `json:"volumes"`
}

type NamespaceWithUserPermissions struct {
	NamespaceWithPermission
	Users []PermissionRecord `json:"users"`
}

type VolumeWithUserPermissions struct {
	VolumeWithPermission
	Users []PermissionRecord `json:"users"`
}
