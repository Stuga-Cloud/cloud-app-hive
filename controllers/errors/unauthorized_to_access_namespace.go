package errors

import "fmt"

type UnauthorizedToAccessNamespaceError struct {
	NamespaceID   string
	NamespaceName string
	UserID        string
}

func (e *UnauthorizedToAccessNamespaceError) Error() string {
	return fmt.Sprintf("user '%s' is unauthorized to access namespace '%s' (namespaceID: %s)", e.UserID, e.NamespaceName, e.NamespaceID)
}

func NewUnauthorizedToAccessNamespaceError(
	namespaceID string,
	namespaceName string,
	userID string,
) *UnauthorizedToAccessNamespaceError {
	return &UnauthorizedToAccessNamespaceError{
		NamespaceID:   namespaceID,
		NamespaceName: namespaceName,
		UserID:        userID,
	}
}
