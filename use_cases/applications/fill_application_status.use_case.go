package applications

import (
	"cloud-app-hive/domain"
	"fmt"

	"cloud-app-hive/domain/commands"
	"cloud-app-hive/domain/repositories"
)

type FillApplicationStatusUseCase struct {
	ContainerManagerRepository repositories.ContainerManagerRepository
}

func (fillApplicationStatusUseCase FillApplicationStatusUseCase) Execute(namespaceName string, applications []domain.Application) ([]domain.Application, error) {
	// Launch a goroutine for each application that calls the ContainerManagerRepository to get the status of the application
	var applicationsWithStatus []domain.Application
	for _, application := range applications {
		status, err := fillApplicationStatusUseCase.ContainerManagerRepository.GetApplicationStatus(commands.GetApplicationStatus{
			Name:      application.Name,
			Namespace: namespaceName,
		})
		if err != nil {
			return nil, fmt.Errorf("error while getting application status: %v", err)
		}
		application.Status, err = domain.KubernetesDeploymentTypeToStatus(status.DeploymentCondition[0].Type)
		if err != nil {
			return nil, fmt.Errorf("error while converting deployment type to status: %v", err)
		}
		applicationsWithStatus = append(applicationsWithStatus, application)
	}

	fmt.Println("Applications with status: ", applicationsWithStatus, len(applicationsWithStatus))
	return applicationsWithStatus, nil
}
