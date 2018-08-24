package model

const ImportSuccessfulMessage = "Successfully imported"

// ImportResponseTotal -- import response for all services
//
// swagger:model
type ImportResponseTotal map[string]ImportResponse

// ImportResponse -- response after resources import
//
// swagger:model
type ImportResponse struct {
	Imported []ImportResult `json:"imported" yaml:"imported"`
	Failed   []ImportResult `json:"failed" yaml:"failed"`
}

// ImportResult -- import result for one resource
//
// swagger:model
type ImportResult struct {
	Name      string `json:"name" yaml:"name"`
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Message   string `json:"message" yaml:"message"`
}

func (resp *ImportResponse) ImportSuccessful(name, namespace string) {
	resp.Imported = append(resp.Imported, ImportResult{
		Name:      name,
		Namespace: namespace,
		Message:   ImportSuccessfulMessage,
	})
}

func (resp *ImportResponse) ImportFailed(name, namespace, message string) {
	resp.Failed = append(resp.Failed, ImportResult{
		Name:      name,
		Namespace: namespace,
		Message:   message,
	})
}
