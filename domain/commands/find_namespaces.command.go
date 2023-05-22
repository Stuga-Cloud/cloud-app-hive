package commands

// FindNamespaces is a command that represents the data needed to find namespaces
type FindNamespaces struct {
	Name    string
	UserID  string
	Page    int
	PerPage int
}
