package applications

import "cloud-app-hive/domain"

type CreateApplicationUseCase struct {
	// All the repositories that the use case needs
}

func (createApplicationUseCase CreateApplicationUseCase) Execute(deployApplication domain.DeployApplication) (domain.Application, error) {
	// All the logic of the use case

	// Get user namespaces and applications

	// Check if the namespace exists

	// Check if the user can access to the namespace

	// Check if the application already exists - by application name and namespace
	return domain.Application{}, nil
}
