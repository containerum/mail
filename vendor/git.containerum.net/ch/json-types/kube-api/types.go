package kube_api

type Protocol string
const (
	UDP Protocol = "UDP"
	TCP Protocol = "TCP"
)

type Endpoint struct {
	Name      string   `json:"name"`
	Owner     *string  `json:"owner,omitempty"`
	CreatedAt *int64   `json:"created_at,omitempty"`
	Addresses []string `json:"addresses"`
	Ports     []Port   `json:"ports"`
}

type Port struct {
	Name       string   `json:"name"`
	Port       int      `json:"port"`
	Protocol   Protocol `json:"protocol"`
}
