package resource

import (
	"time"

	"git.containerum.net/ch/json-types/misc"
	"github.com/lib/pq"
)

type Kind string // constants KindNamespace, KindVolume, ... It`s recommended to use strings.ToLower before comparsion

const (
	KindNamespace  Kind = "namespace"
	KindVolume     Kind = "volume"
	KindExtService Kind = "extservice"
	KindIntService Kind = "intservice"
)

type PermissionStatus string // constants PermissionStatusOwner, PermissionStatusRead

const (
	PermissionStatusOwner      PermissionStatus = "owner"
	PermissionStatusRead       PermissionStatus = "read"
	PermissionStatusWrite      PermissionStatus = "write"
	PermissionStatusReadDelete PermissionStatus = "readdelete"
	PermissionStatusNone       PermissionStatus = "none"
)

type Resource struct {
	ID         string        `json:"id,omitempty" db:"id"`
	CreateTime time.Time     `json:"create_time,omitempty" db:"create_time"`
	Deleted    bool          `json:"deleted,omitempty" db:"deleted"` // not optional because we actually don`t need it if it`s false
	DeleteTime misc.NullTime `json:"delete_time,omitempty" db:"delete_time"`
	TariffID   string        `json:"tariff_id,omitempty" db:"tariff_id"`
}

func (r *Resource) Mask() {
	r.ID = ""
	r.CreateTime = time.Time{}
	r.Deleted = false
	r.DeleteTime.Valid = false
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

	Active      misc.NullBool   `json:"active,omitempty" db:"active"`
	Capacity    int             `json:"capacity" db:"capacity"` // gigabytes
	Replicas    int             `json:"replicas,omitempty" db:"replicas"`
	NamespaceID misc.NullString `json:"namespace_id,omitempty" db:"ns_id"`

	GlusterName string `json:"gluster_name,omitempty" db:"gluster_name"`
	StorageID   string `json:"storage_id,omitempty" db:"storage_id"`
}

func (v *Volume) Mask() {
	v.Resource.Mask()
	v.Active.Valid = false
	v.Replicas = 0
	v.NamespaceID.Valid = false
	v.GlusterName = ""
	v.StorageID = ""
}

type Storage struct {
	ID       string         `json:"id,omitempty" db:"id"`
	Name     string         `json:"name" db:"name"`
	Used     int            `json:"used" db:"used"`
	Size     int            `json:"size" db:"size"`
	Replicas int            `json:"replicas"`
	IPs      pq.StringArray `json:"ips" db:"ips"`
}

type Deployment struct {
	ID          string        `json:"id,omitempty" db:"id"`
	NamespaceID string        `json:"namespace_id,omitempty" db:"ns_id"`
	Name        string        `json:"name" db:"name"`
	CreateTime  time.Time     `json:"create_time,omitempty" db:"create_time"`
	Deleted     bool          `json:"deleted,omitempty" db:"deleted"`
	DeleteTime  misc.NullTime `json:"delete_time,omitempty" db:"delete_time"`
	Replicas    int           `json:"replicas" db:"replicas"`
}

func (d *Deployment) Mask() {
	d.ID = ""
	d.NamespaceID = ""
	d.CreateTime = time.Time{}
	d.Deleted = false
	d.DeleteTime.Valid = false
}

type PermissionRecord struct {
	PermID                string           `json:"perm_id,omitempty" db:"perm_id"`
	Kind                  Kind             `json:"kind,omitempty" db:"kind"`
	ResourceID            misc.NullString  `json:"resource_id,omitempty" db:"resource_id"` // it can be null for resources without tables
	ResourceLabel         string           `json:"label,omitempty" db:"resource_label"`
	OwnerUserID           string           `json:"owner_user_id,omitempty" db:"owner_user_id"`
	CreateTime            time.Time        `json:"create_time,omitempty" db:"create_time"`
	UserID                string           `json:"user_id" db:"user_id"`
	AccessLevel           PermissionStatus `json:"access" db:"access_level"`
	Limited               bool             `json:"limited,omitempty" db:"limited"`
	AccessLevelChangeTime time.Time        `json:"access_level_change_time" db:"access_level_change_time"`
	NewAccessLevel        PermissionStatus `json:"new_access_level,omitempty" db:"new_access_level"`
}

func (p *PermissionRecord) Mask() {
	p.PermID = ""
	p.Kind = "" // will be already known though
	p.ResourceID.Valid = false
	p.OwnerUserID = ""
	p.CreateTime = time.Time{}
	p.UserID = ""
	p.AccessLevel = p.NewAccessLevel
	p.Limited = false
	p.AccessLevelChangeTime = time.Time{}
	p.NewAccessLevel = ""
}

type Container struct {
	ID       string `json:"id,omitempty" db:"id"`
	DeployID string `json:"depl_id,omitempty" db:"depl_id"`
	Name     string `json:"name" db:"name"`
	Image    string `json:"image" db:"image"`
	RAM      int    `json:"ram" db:"ram"`
	CPU      int    `json:"cpu" db:"cpu"`
}

func (c *Container) Mask() {
	c.ID = ""
	c.DeployID = ""
}

type EnvironmentVariable struct {
	EnvID       string `json:"id,omitempty" db:"env_id"`
	ContainerID string `json:"container_id,omitempty" db:"container_id"`
	Name        string `json:"name" db:"name"`
	Value       string `json:"value" db:"value"`
}

func (e *EnvironmentVariable) Mask() {
	e.EnvID = ""
	e.ContainerID = ""
}

type VolumeMount struct {
	MountID     string          `json:"id,omitempty" db:"mount_id"`
	ContainerID string          `json:"container_id,omitempty" db:"container_id"`
	VolumeID    string          `json:"volume_id,omitempty" db:"volume_id"`
	MountPath   string          `json:"mount_path" db:"mount_path"`
	SubPath     misc.NullString `json:"sub_path,omitempty" db:"sub_path"`
}

func (vm *VolumeMount) Mask() {
	vm.MountID = ""
	vm.ContainerID = ""
	vm.VolumeID = ""
}

type Domain struct {
	ID          string         `json:"id,omitempty" binding:"-" db:"id"`
	Domain      string         `json:"domain" binding:"required" db:"domain"`
	DomainGroup string         `json:"domain_group" db:"domain_group"`
	IP          pq.StringArray `json:"ip" binding:"required,dive,ip"`
}

type IngressType string

const (
	IngressHTTP        IngressType = "http"
	IngressHTTPS                   = "https"
	IngressCustomHTTPS             = "custom_https"
)

type IngressEntry struct {
	ID          string      `json:"id,omitempty" db:"id"`
	Domain      string      `json:"domain" db:"custom_domain"`
	Type        IngressType `json:"type" db:"type"`
	ServiceID   string      `json:"service_id" db:"service_id"`
	CreatedAt   time.Time   `json:"created_at" db:"created_at"`
	Path        string      `json:"path" db:"path"`
	ServicePort int         `json:"service_port" db:"service_port"`
}

type ServiceType string

const (
	ServiceInternal ServiceType = "internal"
	ServiceExternal ServiceType = "external"
)

type Service struct {
	ID         string        `json:"id,omitempty" db:"id"`
	DeployID   string        `json:"deployment_id,omitempty" db:"depl_id"`
	Name       string        `json:"name" db:"name"`
	Type       ServiceType   `json:"type" db:"type"`
	CreatedAt  time.Time     `json:"created_at,omitempty" db:"created_at"`
	Deleted    bool          `json:"deleted,omitempty" db:"deleted"`
	DeleteTime misc.NullTime `json:"delete_time,omitempty" db:"delete_time"`
}

func (s *Service) Mask() {
	s.ID = ""
	s.DeployID = ""
	s.CreatedAt = time.Time{}
	s.Deleted = false
	s.DeleteTime.Valid = false
}

type PortProtocol string

const (
	ProtocolTCP PortProtocol = "tcp"
	ProtocolUDP PortProtocol = "udp"
)

type Port struct {
	ID         string       `json:"id,omitempty" db:"id"`
	ServiceID  string       `json:"service_id" db:"service_id"`
	Name       string       `json:"name" db:"name"`
	Port       int          `json:"port" db:"port"`
	TargetPort *int         `json:"target_port" db:"target_port"`
	Protocol   PortProtocol `json:"protocol" db:"protocol"`
	Domain     *string      `json:"domain" db:"domain"`
}

func (p *Port) Mask() {
	p.ID = ""
	p.ServiceID = ""
}

// Types below is not for storing in db

type NamespaceWithPermission struct {
	Namespace
	PermissionRecord
}

func (np *NamespaceWithPermission) Mask() {
	np.Namespace.Mask()
	np.PermissionRecord.Mask()
}

type VolumeWithPermission struct {
	Volume
	PermissionRecord
}

func (vp *VolumeWithPermission) Mask() {
	vp.Volume.Mask()
	vp.PermissionRecord.Mask()
}

type NamespaceWithVolumes struct {
	NamespaceWithPermission
	Volume []VolumeWithPermission `json:"volumes"`
}

func (nv *NamespaceWithVolumes) Mask() {
	nv.NamespaceWithPermission.Mask()
	for i := range nv.Volume {
		nv.Volume[i].Mask()
	}
}

type NamespaceWithUserPermissions struct {
	NamespaceWithPermission
	Users []PermissionRecord `json:"users,omitempty"`
}

func (nu *NamespaceWithUserPermissions) Mask() {
	borrowed := nu.UserID != nu.OwnerUserID
	nu.NamespaceWithPermission.Mask()
	if borrowed {
		nu.Users = nil
	} else {
		for i := range nu.Users {
			nu.Users[i].Mask()
		}
	}
}

type VolumeWithUserPermissions struct {
	VolumeWithPermission
	Users []PermissionRecord `json:"users,omitempty"`
}

func (vp *VolumeWithUserPermissions) Mask() {
	borrowed := vp.UserID != vp.OwnerUserID
	vp.VolumeWithPermission.Mask()
	if borrowed {
		vp.Users = nil
	} else {
		for i := range vp.Users {
			vp.Users[i].Mask()
		}
	}
}
