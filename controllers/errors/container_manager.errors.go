package errors

import (
	"fmt"
	"time"
)

type ContainerManagerError struct {
	Message string
}

func (e *ContainerManagerError) Error() string {
	return e.Message
}

type ContainerManagerConnectionError struct {
	Message string
}

func (e *ContainerManagerConnectionError) Error() string {
	return e.Message
}

type ContainerManagerApplicationDeploymentError struct {
	Message         string
	ApplicationName string
	Namespace       string
	Image           string
}

func (e *ContainerManagerApplicationDeploymentError) Error() string {
	return fmt.Sprintf("%s (application %s in namespace %s with image %s failed to deploy at %s)", e.Message, e.ApplicationName, e.Namespace, e.Image, time.Now().UTC().Format(time.RFC3339))
}

type ContainerManagerApplicationRemoveError struct {
	Message         string
	ApplicationName string
	Namespace       string
}

func (e *ContainerManagerApplicationRemoveError) Error() string {
	return fmt.Sprintf(
		"%s (application %s in namespace %s failed to remove at %s)", e.Message, e.ApplicationName, e.Namespace, time.Now().UTC().Format(time.RFC3339),
	)
}

type ContainerManagerApplicationInformationError struct {
	Message         string
	ApplicationName string
	Namespace       string
	Type            string
}

func (e *ContainerManagerApplicationInformationError) Error() string {
	return fmt.Sprintf("%s (application %s in namespace %s with type %s failed to get information at %s)", e.Message, e.ApplicationName, e.Namespace, e.Type, time.Now().UTC().Format(time.RFC3339))
}

type ContainerManagerNamespaceRemoveError struct {
	Message   string
	Namespace string
}

func (e *ContainerManagerNamespaceRemoveError) Error() string {
	return fmt.Sprintf("%s (namespace %s failed to remove at %s)", e.Message, e.Namespace, time.Now().UTC().Format(time.RFC3339))
}
