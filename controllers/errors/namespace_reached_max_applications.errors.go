package errors

type NamespaceHasReachedMaxNumberOfApplicationsError struct {
	Message string
}

func (e *NamespaceHasReachedMaxNumberOfApplicationsError) Error() string {
	return e.Message
}

func NewNamespaceHasReachedMaxNumberOfApplicationsError(
	message string,
) *NamespaceHasReachedMaxNumberOfApplicationsError {
	return &NamespaceHasReachedMaxNumberOfApplicationsError{
		Message: message,
	}
}
