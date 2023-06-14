package errors

type InvalidApplicationScalabilitySpecificationsError struct {
	Message string
}

func (e *InvalidApplicationScalabilitySpecificationsError) Error() string {
	return e.Message
}

func NewInvalidApplicationScalabilitySpecificationsError(
	message string,
) *InvalidApplicationScalabilitySpecificationsError {
	return &InvalidApplicationScalabilitySpecificationsError{
		Message: message,
	}
}

type LimitUnitScanError struct {
	Message string
}

func (e *LimitUnitScanError) Error() string {
	return e.Message
}

func NewLimitUnitScanError(message string) *LimitUnitScanError {
	return &LimitUnitScanError{
		Message: message,
	}
}

type ContainerLimitScanError struct {
	Message string
}

func (e *ContainerLimitScanError) Error() string {
	return e.Message
}

func NewContainerLimitScanError(message string) *ContainerLimitScanError {
	return &ContainerLimitScanError{
		Message: message,
	}
}

type ContainerLimitValueError struct {
	Message string
}

func (e *ContainerLimitValueError) Error() string {
	return e.Message
}

func NewContainerLimitValueError(message string) *ContainerLimitValueError {
	return &ContainerLimitValueError{
		Message: message,
	}
}

type ApplicationContainerSpecificationsScanError struct {
	Message string
}

func (e *ApplicationContainerSpecificationsScanError) Error() string {
	return e.Message
}

func NewApplicationContainerSpecificationsScanError(message string) *ApplicationContainerSpecificationsScanError {
	return &ApplicationContainerSpecificationsScanError{
		Message: message,
	}
}

type ApplicationContainerSpecificationsValueError struct {
	Message string
}

func (e *ApplicationContainerSpecificationsValueError) Error() string {
	return e.Message
}

func NewApplicationContainerSpecificationsValueError(message string) *ApplicationContainerSpecificationsValueError {
	return &ApplicationContainerSpecificationsValueError{
		Message: message,
	}
}
