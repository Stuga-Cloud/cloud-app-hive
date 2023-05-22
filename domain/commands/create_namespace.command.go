package commands

// CreateNamespace is a command that represents the deployment of a namespace
type CreateNamespace struct {
	Name        string
	Description string
	UserID      string
}
