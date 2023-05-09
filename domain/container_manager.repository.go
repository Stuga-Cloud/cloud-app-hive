package domain

// ContainerManagerRepository is an interface that represents a repository of container managers
type ContainerManagerRepository interface {
	// GetMetricsOfApplication returns the metrics of an application
	GetMetricsOfApplication(namespace, applicationName string) ([]ApplicationMetrics, error)
	// DeployApplication deploys an application on a container manager
	DeployApplication(deployApplication DeployApplication) (string, error)
}
