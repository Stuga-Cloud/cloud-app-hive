package schedulers

import (
	"cloud-app-hive/database"
	"cloud-app-hive/repositories"
	"cloud-app-hive/services"
	"cloud-app-hive/use_cases/applications"
)

func InitSchedulers() {
	db, err := database.ConnectToDatabase()
	if err != nil {
		panic(err)
	}

	if err = database.MigrateDatabase(db); err != nil {
		panic(err)
	}

	applicationRepository := repositories.GORMApplicationRepository{
		Database: db,
	}
	containerManager := repositories.KubernetesContainerManagerRepository{}

	retrieveManualScalingApplicationsUseCase := applications.FindManualScalingApplicationsUseCase{
		ApplicationRepository: applicationRepository,
	}
	getApplicationMetricsUseCase := applications.GetApplicationMetricsUseCase{
		ContainerManagerRepository: containerManager,
	}
	emailService := services.NewEmailService()
	scalabilityNotificationService := services.NewScalabilityNotificationService(*emailService)

	notifyApplicationScalingRecommendationScheduler(retrieveManualScalingApplicationsUseCase, getApplicationMetricsUseCase, *scalabilityNotificationService)

	AutoScaleApplicationsAndNotify()
}
