package models

// Service model to expose the information
// about application running on cluster.
type Service struct {
	// the deployment of Name
	Name string `json:"name,omitempty"`
	// the deployment belongs to which ApplicationGroup label
	ApplicationGroup string `json:"applicationGroup,omitempty"`
	// total number of running pod corresponding to serviceName
	RunningPodsCount int `json:"runningPodsCount,omitempty"`
}
