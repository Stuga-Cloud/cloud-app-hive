package errors

import "fmt"

type NamespaceNotFoundByNameError struct {
	NamespaceName string
}

func (e *NamespaceNotFoundByNameError) Error() string {
	return fmt.Sprintf("namespace with name %s not found", e.NamespaceName)
}

type NamespaceNotFoundByIDError struct {
	NamespaceID string
}

func (e *NamespaceNotFoundByIDError) Error() string {
	return fmt.Sprintf("namespace with id %s not found", e.NamespaceID)
}

func NewNamespaceNotFoundByNameError(
	namespaceName string,
) *NamespaceNotFoundByNameError {
	return &NamespaceNotFoundByNameError{
		NamespaceName: namespaceName,
	}
}

func NewNamespaceNotFoundByIDError(
	namespaceID string,
) *NamespaceNotFoundByIDError {
	return &NamespaceNotFoundByIDError{
		NamespaceID: namespaceID,
	}
}
