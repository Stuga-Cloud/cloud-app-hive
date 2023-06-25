package commands

import "cloud-app-hive/domain"

// UpdateApplication is a command that represents the deployment of an application
type UpdateApplication struct {
	UserID                    string
	Description               string
	Image                     string
	Registry                  domain.ImageRegistry
	Port                      uint32
	ApplicationType           domain.ApplicationType
	EnvironmentVariables      domain.ApplicationEnvironmentVariables
	Secrets                   domain.ApplicationSecrets
	ContainerSpecifications   domain.ApplicationContainerSpecifications
	ScalabilitySpecifications domain.ApplicationScalabilitySpecifications
	AdministratorEmail        string
}
