package model

// SolutionsTemplatesList -- list of available solutions
//
// swagger:model
type SolutionsTemplatesList struct {
	Solutions []SolutionTemplate `json:"solutions"`
}

func (list SolutionsTemplatesList) Len() int {
	return len(list.Solutions)
}

func (list SolutionsTemplatesList) Copy() SolutionsTemplatesList {
	var solutions = make([]SolutionTemplate, 0, list.Len())
	for _, sol := range solutions {
		solutions = append(solutions, sol.Copy())
	}
	return SolutionsTemplatesList{
		Solutions: solutions,
	}
}

func (list SolutionsTemplatesList) Get(i int) SolutionTemplate {
	return list.Solutions[i]
}

func (list SolutionsTemplatesList) Filter(pred func(SolutionTemplate) bool) SolutionsTemplatesList {
	solutions := make([]SolutionTemplate, 0, list.Len())
	for _, sol := range list.Solutions {
		if pred(sol.Copy()) {
			solutions = append(solutions, sol.Copy())
		}
	}
	return SolutionsTemplatesList{
		Solutions: solutions,
	}
}

// SolutionTemplate -- solution which user can run
//
// swagger:model
type SolutionTemplate struct {
	ID     string          `json:"id,omitempty" yaml:"id,omitempty"`
	Name   string          `json:"name" yaml:"name"`
	Limits *SolutionLimits `json:"limits" yaml:"limits"`
	Images []string        `json:"images" yaml:"images"`
	URL    string          `json:"url" yaml:"url"`
	Active bool            `json:"active" yaml:"active"`
}

func (solution SolutionTemplate) Copy() SolutionTemplate {
	return SolutionTemplate{
		Name: solution.Name,
		Limits: func() *SolutionLimits {
			if solution.Limits == nil {
				return nil
			}
			cp := *solution.Limits
			return &cp
		}(),
		Images: append([]string{}, solution.Images...),
		URL:    solution.URL,
		Active: solution.Active,
	}
}

// SolutionLimits -- solution resources limits
//
// swagger:model
type SolutionLimits struct {
	CPU uint `json:"cpu" yaml:"cpu"`
	RAM uint `json:"ram" yaml:"ram"`
}

// SolutionEnv -- solution environment variables
//
// swagger:model
type SolutionEnv struct {
	Env map[string]string `json:"env"`
}

// SolutionResources -- list of solution resources
//
// swagger:model
type SolutionResources struct {
	Resources map[string]int `json:"resources"`
}

func (res SolutionResources) Copy() SolutionResources {
	r := make(map[string]int, len(res.Resources))
	for k, v := range res.Resources {
		r[k] = v
	}
	return SolutionResources{
		Resources: r,
	}
}

type ConfigFile struct {
	Name string `json:"config_file"`
	Type string `json:"type"`
}

// SolutionsList -- list of running solution
//
// swagger:model
type SolutionsList struct {
	Solutions []Solution `json:"solutions"`
}

func (list SolutionsList) Copy() SolutionsList {
	var solutions = make([]Solution, 0, list.Len())
	for _, s := range solutions {
		solutions = append(solutions, s.Copy())
	}
	return SolutionsList{
		Solutions: solutions,
	}
}

func (list SolutionsList) Len() int {
	return len(list.Solutions)
}

func (list SolutionsList) Get(i int) Solution {
	return list.Solutions[i]
}

func (list SolutionsList) Filter(pred func(Solution) bool) SolutionsList {
	solutions := make([]Solution, 0, list.Len())
	for _, sol := range list.Solutions {
		if pred(sol.Copy()) {
			solutions = append(solutions, sol.Copy())
		}
	}
	return SolutionsList{
		Solutions: solutions,
	}
}

// Solution -- running solution
//
// swagger:model
type Solution struct {
	ID     string            `json:"id,omitempty"`
	Branch string            `json:"branch"`
	Env    map[string]string `json:"env"`
	URL    string            `json:"url"`
	// required: true
	Template string `json:"template"`
	// required: true
	Name string `json:"name"`
	// required: true
	Namespace string `json:"namespace"`
}

func (solution Solution) Copy() Solution {
	env := make(map[string]string, len(solution.Env))
	for k, v := range solution.Env {
		env[k] = v
	}
	return Solution{
		Branch:    solution.Branch,
		Env:       env,
		Template:  solution.Template,
		Name:      solution.Name,
		Namespace: solution.Namespace,
	}
}

// RunSolutionResponse -- response to run solution request
//
// swagger:model
type RunSolutionResponse struct {
	Created    int      `json:"created"`
	NotCreated int      `json:"not_created"`
	Errors     []string `json:"errors,omitempty"`
}
