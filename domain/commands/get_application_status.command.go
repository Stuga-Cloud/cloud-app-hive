package commands

// GetApplicationStatus is a command that represents a request to get the metrics of an application
type GetApplicationStatus struct {
	Name      string
	Namespace string
}
