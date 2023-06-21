package schedulers

import (
	"cloud-app-hive/domain"
	"cloud-app-hive/domain/commands"
	"cloud-app-hive/services"
	"cloud-app-hive/use_cases/applications"
	"cloud-app-hive/utils"
	"fmt"
	"os"
	"strconv"
	"time"
)

func notifyApplicationScalingRecommendationScheduler(
	findManualScalingApplicationsUseCase applications.FindManualScalingApplicationsUseCase,
	getApplicationMetricsUseCase applications.GetApplicationMetricsUseCase,
	scalabilityNotificationService services.ScalabilityNotificationService,
) {
	fmt.Println("Starting NotifyApplicationScalingRecommendation scheduler...")
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered from panic during NotifyApplicationScalingRecommendationScheduler:", r)
			}
		}()
		repeatInterval := getSchedulerRepeatInterval()
		ticker := time.NewTicker(time.Duration(repeatInterval) * time.Second)
		for {
			select {
			case <-ticker.C:
				fmt.Println("=====================================================")
				// TODO
				// 1. get all applications that are manual scaling
				foundApplications, err := findManualScalingApplicationsUseCase.Execute()
				if err != nil {
					fmt.Println("error when try to get applications :", err.Error())
					panic("Error when try to get applications during NotifyApplicationScalingRecommendationScheduler : " + err.Error())
				}
				if len(foundApplications) == 0 {
					fmt.Println("No applications found")
				}

				// 2. parallelize checking each application through the container manager with goroutines
				routines := len(foundApplications)
				done := make(chan bool, routines)
				for _, application := range foundApplications {
					go func(application domain.Application) {
						metrics, err := getApplicationMetricsUseCase.Execute(commands.GetApplicationMetrics{
							Name:      application.Name,
							Namespace: application.Namespace.Name,
						})
						if err != nil {
							fmt.Println("error when try to get application metrics during NotifyApplicationScalingRecommendationScheduler :", err.Error())
							done <- true
							return
						}
						if len(metrics) == 0 {
							fmt.Println("No metrics found for application", application.Name)
							done <- true
							return
						}

						// for _, metric := range metrics {
						// 	fmt.Println("Metrics of application", application.Name, ":", metric.String())
						// }

						var compareActualUsageToAcceptedPercentageResults []domain.CompareActualUsageToAcceptedPercentageResult
						acceptedCPUUsageThreshold := application.ScalabilitySpecifications.CpuUsagePercentageThreshold
						acceptedMemoryUsageThreshold := application.ScalabilitySpecifications.MemoryUsagePercentageThreshold
						for _, metric := range metrics {
							doesExceedCPUAcceptedPercentage, actualCPUUsage := utils.DoesUsageExceedsLimitAndHowMuchActually(metric.CPUUsage, metric.MaxCPUUsage, acceptedCPUUsageThreshold)
							doesExceedMemoryAcceptedPercentage, actualMemoryUsage := utils.DoesUsageExceedsLimitAndHowMuchActually(metric.MemoryUsage, metric.MaxMemoryUsage, acceptedMemoryUsageThreshold)
							// doesExceedEphemeralStorageAcceptedPercentage, actualEphemeralStorageUsage := utils.DoesUsageExceedsLimitAndHowMuchActually(metric.EphemeralStorageUsage, metric.MaxEphemeralStorage, acceptedUsagePercentage)
							compareActualUsageToAcceptedPercentageResult := domain.CompareActualUsageToAcceptedPercentageResult{
								CPUUsageResult: domain.UsageComparisonResult{
									DoesExceedAcceptedPercentage: doesExceedCPUAcceptedPercentage,
									ActualUsage:                  actualCPUUsage,
								},
								MemoryUsageResult: domain.UsageComparisonResult{
									DoesExceedAcceptedPercentage: doesExceedMemoryAcceptedPercentage,
									ActualUsage:                  actualMemoryUsage,
								},
								// EphemeralStorageUsageResult: domain.UsageComparisonResult{
								// 	DoesExceedAcceptedPercentage: doesExceedEphemeralStorageAcceptedPercentage,
								// 	ActualUsage:                  actualEphemeralStorageUsage,
								// },
								PodName: metric.PodName,
							}
							compareActualUsageToAcceptedPercentageResults = append(compareActualUsageToAcceptedPercentageResults, compareActualUsageToAcceptedPercentageResult)
						}

						for _, compareActualUsageToAcceptedPercentageResult := range compareActualUsageToAcceptedPercentageResults {
							// 3. if the application is using more than 80% of its resources, recommend to scale up
							// 4. if the application is using less than 20% of its resources, recommend to scale down
							// 5. send email notification to the application owner
							if compareActualUsageToAcceptedPercentageResult.OneOfTheUsageExceedsAcceptedPercentage() {
								success, err := scalabilityNotificationService.SendApplicationScalabilityRecommandationMail(
									application.AdministratorEmail,
									application.Name,
									application.Namespace.Name,
									compareActualUsageToAcceptedPercentageResult,
									application.ContainerSpecifications,
									application.ScalabilitySpecifications,
								)
								if err != nil {
									fmt.Println("error when try to send application scalability recommandation mail :", err.Error())
									panic("Error when try to send application scalability recommandation mail during NotifyApplicationScalingRecommendationScheduler : " + err.Error())
								}
								if success {
									fmt.Println("Application", application.Name, "is exceeding accepted percentages, email sent to", application.AdministratorEmail)
								} else {
									fmt.Println("Application", application.Name, "is exceeding accepted percentages, but email not sent to", application.AdministratorEmail)
								}
							} else {
								fmt.Println("Application", application.Name, "is not exceeding accepted percentages")
							}
						}

						done <- true
					}(application)
				}
				for i := 0; i < routines; i++ {
					<-done
				}
				return
			}
		}
	}()
}

func getSchedulerRepeatInterval() int {
	const defaultRepeatInterval = 1
	var repeatInterval int
	SchedulerRecommendApplicationScalingInMinutes := os.Getenv("SCHEDULER_RECOMMEND_APPLICATION_SCALING_IN_SECONDS")
	if SchedulerRecommendApplicationScalingInMinutes == "" {
		fmt.Println("SCHEDULER_RECOMMEND_APPLICATION_SCALING_IN_SECONDS is not set, using default value")
		repeatInterval = defaultRepeatInterval
	}
	repeatInterval, err := strconv.Atoi(SchedulerRecommendApplicationScalingInMinutes)
	if err != nil {
		fmt.Println("Error when convert SchedulerRecommendApplicationScalingInMinutes to int")
		panic(fmt.Sprintf("Error when convert SchedulerRecommendApplicationScalingInMinutes to int during NotifyApplicationScalingRecommendationScheduler: %s", err.Error()))
	}
	return repeatInterval
}
