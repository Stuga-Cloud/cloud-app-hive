package commands

// UpdateNamespace is a command that represents a user's request to update a namespace
type UpdateNamespace struct {
	ID          string
	Description string
	UserID      string
}
