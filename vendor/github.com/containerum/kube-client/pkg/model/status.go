package model

//swagger:model
type ServiceStatus struct {
	Name     string            `json:"name"`
	Version  string            `json:"version"`
	StatusOK bool              `json:"ok"`
	Details  map[string]string `json:"details,omitempty"`
}

func (status *ServiceStatus) AddDetails(key, value string) *ServiceStatus {
	status.Details[key] = value
	return status
}
