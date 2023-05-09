package domain

// ContainerManagerRepository is an interface that represents a repository of container managers
type ContainerManagerRepository interface {
	// GetApplicationMetrics returns the metrics of an application
	GetApplicationMetrics(application GetApplicationMetrics) ([]ApplicationMetrics, error)
	// DeployApplication deploys an application on a container manager
	DeployApplication(deployApplication DeployApplication) error
	// GetApplicationLogs returns the logs of an application
	GetApplicationLogs(application GetApplicationLogs) ([]ApplicationLogs, error)
}
