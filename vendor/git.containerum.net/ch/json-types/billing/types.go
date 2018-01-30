package billing

import (
	"time"
)

type Tariff struct {
	ID          string    `json:"tariff_id"`
	Label       string    `json:"label"`
	Price       float64   `json:"price"` // TODO: use special fixed-precision type like "golang.org/x/text/currency".Amount
	Active      bool      `json:"is_active"`
	Public      bool      `json:"is_public"`
	BillingID   int       `json:"billing_id"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type NamespaceTariff struct {
	Tariff

	CPULimit         int     `json:"cpu_limit"`
	MemoryLimit      int     `json:"memory_limit"`
	Traffic          int     `json:"traffic"`       // in gigabytes
	TrafficPrice     float64 `json:"traffic_price"` // price of overused traffic. TODO: @see Tariff.Price
	ExternalServices int     `json:"external_services"`
	InternalServices int     `json:"internal_services"`
}

type VolumeTariff struct {
	Tariff

	StorageLimit  int `json:"storage_limit"` // in gigabytes
	ReplicasLimit int `json:"replicas_limit"`
}
