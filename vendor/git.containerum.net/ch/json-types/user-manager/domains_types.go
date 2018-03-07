package user

import "time"

type DomainListResponse struct {
	DomainList []Domain `json:"domain_list,omitempty"`
}

type Domain struct {
	Domain    string    `json:"domain"`
	AddedBy   string    `json:"added_by,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
