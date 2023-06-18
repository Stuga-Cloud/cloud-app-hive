package errors

import "fmt"

type ApplicationNotFoundByIDError struct {
	ApplicationID string
}

func (e *ApplicationNotFoundByIDError) Error() string {
	return fmt.Sprintf("application with id %s not found", e.ApplicationID)
}

func NewApplicationNotFoundByIDError(
	ApplicationID string,
) *ApplicationNotFoundByIDError {
	return &ApplicationNotFoundByIDError{
		ApplicationID: ApplicationID,
	}
}
