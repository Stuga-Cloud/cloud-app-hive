package commands

import "cloud-app-hive/domain"

// CreateApplication is a command that represents the deployment of an application
type CreateApplication struct {
	Name                      string
	Description               string
	Image                     string
	NamespaceID               string
	UserID                    string
	Port                      uint32
	Zone                      string
	ApplicationType           domain.ApplicationType
	EnvironmentVariables      domain.ApplicationEnvironmentVariables
	Secrets                   domain.ApplicationSecrets
	ContainerSpecifications   domain.ApplicationContainerSpecifications
	ScalabilitySpecifications domain.ApplicationScalabilitySpecifications
	AdministratorEmail        string
}

// ToDomain converts the command to a domain.Application
func (createApplication CreateApplication) ToDomain(ID string) domain.Application {
	return domain.Application{
		ID:                        ID,
		Name:                      createApplication.Name,
		Description:               createApplication.Description,
		Image:                     createApplication.Image,
		NamespaceID:               createApplication.NamespaceID,
		UserID:                    createApplication.UserID,
		Port:                      createApplication.Port,
		Zone:                      createApplication.Zone,
		ApplicationType:           createApplication.ApplicationType,
		EnvironmentVariables:      &createApplication.EnvironmentVariables,
		Secrets:                   &createApplication.Secrets,
		ContainerSpecifications:   &createApplication.ContainerSpecifications,
		ScalabilitySpecifications: &createApplication.ScalabilitySpecifications,
	}
}
