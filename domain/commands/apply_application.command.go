package commands

import "cloud-app-hive/domain"

// ApplyApplication is a command that represents the deployment of an application
type ApplyApplication struct {
	Name                      string
	Image                     string
	Namespace                 string
	Port                      uint32
	ApplicationType           domain.ApplicationType
	EnvironmentVariables      domain.ApplicationEnvironmentVariables
	Secrets                   domain.ApplicationSecrets
	ContainerSpecifications   domain.ApplicationContainerSpecifications
	ScalabilitySpecifications domain.ApplicationScalabilitySpecifications
}
