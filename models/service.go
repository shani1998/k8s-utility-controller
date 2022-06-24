package models

// Service to expose data about pods
type Service struct {
	// the pod of Name
	Name string `json:"name"`
	// the pod belongs to which ApplicationGroup label
	ApplicationGroup string `json:"applicationGroup"`
	// total number of running pod in given configuration
	RunningPodsCount int `json:"runningPodsCount"`
}
