package responses

import "cloud-app-hive/domain"

// CreateApplicationResponse is a struct that represents the response body for creating an application
// swagger:response CreateApplicationResponse
type CreateApplicationResponse struct {
	Message     string              `json:"message"`
	Application ApplicationResponse `json:"application"`
}

// ApplicationResponse is a struct that represents the response body for creating an application
// swagger:model ApplicationResponse
type ApplicationResponse struct {
	ID                        string                                       `json:"id"`
	Name                      string                                       `json:"name"`
	Image                     string                                       `json:"image"`
	NamespaceID               string                                       `json:"namespace_id"`
	UserID                    string                                       `json:"user_id"`
	Port                      uint32                                       `json:"port"`
	ApplicationType           domain.ApplicationType                       `json:"application_type"`
	EnvironmentVariables      *domain.ApplicationEnvironmentVariables      `json:"environment_variables"`
	Secrets                   *domain.ApplicationSecrets                   `json:"secrets"`
	ContainerSpecifications   *domain.ApplicationContainerSpecifications   `json:"container_specifications"`
	ScalabilitySpecifications *domain.ApplicationScalabilitySpecifications `json:"scalability_specifications"`
}

// ApplicationDomainToResponse is a method that converts a domain.Application to an ApplicationResponse
func ApplicationDomainToResponse(application *domain.Application) ApplicationResponse {
	return ApplicationResponse{
		ID:                        application.ID,
		Name:                      application.Name,
		Image:                     application.Image,
		NamespaceID:               application.NamespaceID,
		UserID:                    application.UserID,
		Port:                      application.Port,
		ApplicationType:           application.ApplicationType,
		EnvironmentVariables:      application.EnvironmentVariables,
		Secrets:                   application.Secrets,
		ContainerSpecifications:   application.ContainerSpecifications,
		ScalabilitySpecifications: application.ScalabilitySpecifications,
	}
}
