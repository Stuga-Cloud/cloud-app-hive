package repositories

import (
	"cloud-app-hive/domain"
	"cloud-app-hive/domain/commands"
)

// ContainerManagerRepository is an interface that represents a repository of container managers
type ContainerManagerRepository interface {
	// GetApplicationMetrics returns the metrics of an application
	GetApplicationMetrics(application commands.GetApplicationMetrics) ([]domain.ApplicationMetrics, error)
	// ApplyApplication deploys an application on a container manager
	ApplyApplication(applyApplication commands.ApplyApplication) error
	// GetApplicationLogs returns the logs of an application
	GetApplicationLogs(application commands.GetApplicationLogs) ([]domain.ApplicationLogs, error)
	// GetApplicationStatus returns the status of an application
	GetApplicationStatus(application commands.GetApplicationStatus) (*domain.ApplicationStatus, error)
	// UnapplyApplication delete an application on a container manager
	UnapplyApplication(applyApplication commands.UnapplyApplication) error
	// DeleteNamespace deletes a namespace on a container manager
	DeleteNamespace(namespace string) error
}
