package errors

import "fmt"

type UnauthorizedToRemoveAdminFromNamespaceError struct {
	NamespaceID string
	UserID      string
	RemovedBy   string
}

func (e *UnauthorizedToRemoveAdminFromNamespaceError) Error() string {
	return fmt.Sprintf("user '%s' is unauthorized to remove admin %s from namespace '%s'", e.RemovedBy, e.UserID, e.NamespaceID)
}

func NewUnauthorizedToRemoveAdminFromNamespaceError(
	namespaceID string,
	userID string,
	removedBy string,
) *UnauthorizedToRemoveAdminFromNamespaceError {
	return &UnauthorizedToRemoveAdminFromNamespaceError{
		NamespaceID: namespaceID,
		UserID:      userID,
		RemovedBy:   removedBy,
	}
}
