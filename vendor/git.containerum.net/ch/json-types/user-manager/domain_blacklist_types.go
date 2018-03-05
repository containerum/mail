package user

type DomainToBlacklistRequest struct {
	Domain string `json:"domain" binding:"required"`
}

type DomainListResponse struct {
	DomainList []DomainResponse `json:"domain_list"`
}

type DomainResponse struct {
	Domain    string `json:"domain"`
	AddedBy   string `json:"added_by"`
	CreatedAt string `json:"created_at"`
}