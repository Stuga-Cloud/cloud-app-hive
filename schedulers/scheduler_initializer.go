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

	findManualScalingApplicationsUseCase := applications.FindManualScalingApplicationsUseCase{
		ApplicationRepository: applicationRepository,
	}
	findAutoScalingApplicationsUseCase := applications.FindAutoScalingApplicationsUseCase{
		ApplicationRepository: applicationRepository,
	}
	getApplicationMetricsUseCase := applications.GetApplicationMetricsUseCase{
		ContainerManagerRepository: containerManager,
	}
	emailService := services.NewEmailService()
	scalabilityNotificationService := services.NewScalabilityNotificationService(*emailService)

	manualScaleScheduler := NotifyApplicationScalingRecommendationScheduler{
		findManualScalingApplicationsUseCase,
		getApplicationMetricsUseCase,
		*scalabilityNotificationService,
	}
	manualScaleScheduler.Launch()

	scaleApplicationUseCase := applications.ScaleApplicationUseCase{
		ApplicationRepository: applicationRepository,
	}
	autoScaleScheduler := AutoScaleApplicationsAndNotifyScheduler{
		findAutoScalingApplicationsUseCase,
		getApplicationMetricsUseCase,
		*scalabilityNotificationService,
		scaleApplicationUseCase,
	}

	autoScaleScheduler.Launch()
}
