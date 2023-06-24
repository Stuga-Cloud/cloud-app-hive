package responses

import "cloud-app-hive/domain"

// CreateApplicationResponse is a struct that represents the response body for creating an application
// swagger:response CreateApplicationResponse
type CreateApplicationResponse struct {
	Message     string             `json:"message"`
	Application domain.Application `json:"application"`
}

// ApplicationResponse is a struct that represents the response body for creating an application
// swagger:model ApplicationResponse
type ApplicationResponse struct {
	ID                        string                                       `json:"id"`
	Name                      string                                       `json:"name"`
	Image                     string                                       `json:"image"`
	NamespaceID               string                                       `json:"namespaceId"`
	UserID                    string                                       `json:"userId"`
	Port                      uint32                                       `json:"port"`
	ApplicationType           domain.ApplicationType                       `json:"applicationType"`
	EnvironmentVariables      *domain.ApplicationEnvironmentVariables      `json:"environmentVariables"`
	Secrets                   *domain.ApplicationSecrets                   `json:"secrets"`
	ContainerSpecifications   *domain.ApplicationContainerSpecifications   `json:"containerSpecifications"`
	ScalabilitySpecifications *domain.ApplicationScalabilitySpecifications `json:"scalabilitySpecifications"`
}
