package domain

// GetApplicationMetrics is a command that represents a request to get the metrics of an application
type GetApplicationMetrics struct {
	Name      string
	Namespace string
}
