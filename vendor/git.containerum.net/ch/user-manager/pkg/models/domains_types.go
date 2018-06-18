package models

import "time"

// DomainListResponse -- domains list
//
// swagger:model
type DomainListResponse struct {
	DomainList []Domain `json:"domain_list,omitempty"`
}

// DomainListResponse -- domains list
//
// swagger:model
type Domain struct {
	// required: true
	Domain    string    `json:"domain"`
	AddedBy   string    `json:"added_by,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
