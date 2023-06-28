package applications

import (
	"fmt"

	"cloud-app-hive/domain"
	"cloud-app-hive/domain/commands"
	"cloud-app-hive/domain/repositories"
)

type ScaleApplicationUseCase struct {
	ApplicationRepository repositories.ApplicationRepository
}

type ScalingType string

const (
	HorizontalUpScaling   ScalingType = "HorizontalUpScaling"
	HorizontalDownScaling ScalingType = "HorizontalDownScaling"
	VerticalUpScaling     ScalingType = "VerticalUpScaling"
	VerticalDownScaling   ScalingType = "VerticalDownScaling"
)

func (scaleApplicationUseCase ScaleApplicationUseCase) Execute(applicationID string, updateApplication commands.UpdateApplication, scalingType ScalingType) (*domain.Application, error) {
	foundApplicationByID, err := scaleApplicationUseCase.ApplicationRepository.FindByID(applicationID)
	if err != nil {
		return nil, fmt.Errorf("error while finding application by id: %w", err)
	}
	if foundApplicationByID == nil {
		return nil, fmt.Errorf("no application found for application id %s", applicationID)
	}

	updatedApplication := &domain.Application{}
	if scalingType == HorizontalUpScaling {
		updatedApplication, err = scaleApplicationUseCase.ApplicationRepository.HorizontalScaleUp(applicationID)
	}
	if scalingType == HorizontalDownScaling {
		updatedApplication, err = scaleApplicationUseCase.ApplicationRepository.HorizontalScaleDown(applicationID)
	}
	// if scalingType == VerticalUpScaling {
	// 	updatedApplication, err = scaleApplicationUseCase.ApplicationRepository.VerticalScaleUp(applicationID)
	// }
	// if scalingType == VerticalDownScaling {
	// 	updatedApplication, err = scaleApplicationUseCase.ApplicationRepository.VerticalScaleDown(applicationID)
	// }

	if err != nil {
		return nil, err
	}

	return updatedApplication, nil
}
