package applications

import (
	"cloud-app-hive/domain"
	"cloud-app-hive/domain/commands"
)

type DeleteApplicationUseCase struct {
	// All the repositories that the use case needs
}

func (deleteApplicationUseCase DeleteApplicationUseCase) Execute(deleteApplication commands.DeleteApplication) (domain.Application, error) {
	// TODO -> Delete application in database

	// All the logic of the use case

	// Get user namespaces and applications

	// Check if the namespace exists

	// Check if the user can access to the namespace

	// Check if the application already exists - by application name and namespace

	// Delete the application

	return domain.Application{}, nil
}
