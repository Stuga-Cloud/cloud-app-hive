package commands

// RemoveNamespaceMembership is a command that represents the intent to create a namespace membership.
type RemoveNamespaceMembership struct {
	UserID      string
	NamespaceID string
	RemovedBy   string
}
