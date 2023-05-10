package commands

import "cloud-app-hive/domain"

// ApplyApplication is a command that represents the deployment of an application
type ApplyApplication struct {
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
