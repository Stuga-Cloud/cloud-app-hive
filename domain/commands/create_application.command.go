package commands

import (
	"cloud-app-hive/domain"
)

// CreateApplication is a command that represents the deployment of an application
type CreateApplication struct {
	UserID                    string
	Name                      string
	Description               string
	Image                     string
	NamespaceID               string
	Port                      uint32
	Zone                      string
	ApplicationType           domain.ApplicationType
	EnvironmentVariables      domain.ApplicationEnvironmentVariables
	Secrets                   domain.ApplicationSecrets
	ContainerSpecifications   domain.ApplicationContainerSpecifications
	ScalabilitySpecifications domain.ApplicationScalabilitySpecifications
	AdministratorEmail        string
}
