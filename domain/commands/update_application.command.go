package commands

import "cloud-app-hive/domain"

// UpdateApplication is a command that represents the deployment of an application
type UpdateApplication struct {
	Name                      string
	Image                     string
	Namespace                 string
	Port                      int
	ApplicationType           domain.ApplicationType
	EnvironmentVariables      []domain.ApplicationEnvironmentVariable
	Secrets                   []domain.ApplicationSecret
	ContainerSpecifications   domain.ApplicationContainerSpecifications
	ScalabilitySpecifications domain.ApplicationScalabilitySpecifications
}
