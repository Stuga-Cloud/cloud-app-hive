package commands

// UnapplyApplication is a command that represents a request to get the metrics of an application
type UnapplyApplication struct {
	Name      string
	Namespace string
}
