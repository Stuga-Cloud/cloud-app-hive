package applications

import (
	"cloud-app-hive/domain"
	"cloud-app-hive/domain/commands"
)

type UpdateApplicationUseCase struct {
	// All the repositories that the use case needs
}

func (createApplicationUseCase UpdateApplicationUseCase) Execute(updateApplication commands.UpdateApplication) (domain.Application, error) {
	// TODO -> Update application in database

	// All the logic of the use case

	// Get user namespaces and applications

	// Check if the namespace exists

	// Check if the user can access to the namespace

	// Check if the application already exists - by application name and namespace

	// Update the application

	return domain.Application{}, nil
}
