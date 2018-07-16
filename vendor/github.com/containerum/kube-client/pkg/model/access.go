package model

type UserAccess struct {
	Username    string          `json:"username"`
	AccessLevel UserGroupAccess `json:"access_level"`
}

func (access UserAccess) String() string {
	return access.Username + ":" + string(access.AccessLevel)
}

// ResourceUpdateUserAccess -- contains user access data
//swagger:model
type ResourceUpdateUserAccess struct {
	Username string          `json:"username"`
	Access   UserGroupAccess `json:"access,omitempty"`
}
