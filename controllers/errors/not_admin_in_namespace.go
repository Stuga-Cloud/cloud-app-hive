package errors

import "fmt"

type NotAdminInNamespaceError struct {
	NamespaceID string
	UserID      string
}

func (e *NotAdminInNamespaceError) Error() string {
	return fmt.Sprintf("user '%s' is not an admin in namespace '%s' (namespaceID: %s)", e.UserID, e.NamespaceID, e.NamespaceID)
}

func NewNotAdminInNamespaceError(
	namespaceID string,
	userID string,
) *NotAdminInNamespaceError {
	return &NotAdminInNamespaceError{
		NamespaceID: namespaceID,
		UserID:      userID,
	}
}
