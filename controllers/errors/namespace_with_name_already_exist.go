package errors

import "fmt"

type NamespaceWithNameAlreadyExistError struct {
	NamespaceName string
}

func (e *NamespaceWithNameAlreadyExistError) Error() string {
	return fmt.Sprintf("namespace with name %s already exist", e.NamespaceName)
}

func NewNamespaceWithNameAlreadyExistError(
	namespaceName string,
) *NamespaceWithNameAlreadyExistError {
	return &NamespaceWithNameAlreadyExistError{
		NamespaceName: namespaceName,
	}
}
