package domain

// GetApplicationLogs is a command that represents a request to get the logs of an application
type GetApplicationLogs struct {
	Name      string
	Namespace string
}
