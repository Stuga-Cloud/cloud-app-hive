package schedulers

import (
	"cloud-app-hive/domain"
	"cloud-app-hive/domain/commands"
	"cloud-app-hive/services"
	"cloud-app-hive/use_cases/applications"
	"fmt"
	"os"
	"strconv"
	"time"
)

type AutoScaleApplicationsAndNotifyScheduler struct {
	findAutoScalingApplicationsUseCase applications.FindAutoScalingApplicationsUseCase
	getApplicationMetricsUseCase       applications.GetApplicationMetricsUseCase
	scalabilityNotificationService     services.ScalabilityNotificationService
	scaleApplicationUseCase            applications.ScaleApplicationUseCase
}

func (scheduler AutoScaleApplicationsAndNotifyScheduler) Launch() {
	fmt.Println("Starting 'AutoScaleApplicationsAndNotifyScheduler' scheduler...")
	go func() {
		repeatInterval, err := getAutoScaleAppSchedulerRepeatInterval()
		if err != nil {
			fmt.Println("Error when try to get auto scale app scheduler repeat interval :", err.Error())
			return
		}
		ticker := time.NewTicker(time.Duration(repeatInterval) * time.Second)

		lastNotifiedMaxReplicasDatetimeByApplication := make(map[string]time.Time)

		for {
			select {
			case <-ticker.C:
				// 1. get all applications that are auto scaling
				foundApplications, err := scheduler.findAutoScalingApplicationsUseCase.Execute()
				if err != nil {
					fmt.Println("error when try to get auto scaling applications :", err.Error())
					return
				}
				if len(foundApplications) == 0 {
					fmt.Println("No auto scaling applications found")
					continue
				}

				// 2. parallelize checking each application through the container manager with goroutines
				routines := len(foundApplications)
				done := make(chan bool, routines)
				for _, application := range foundApplications {
					go func(application domain.Application) {
						metrics, err := scheduler.getApplicationMetricsUseCase.Execute(commands.GetApplicationMetrics{
							Name:      application.Name,
							Namespace: application.Namespace.Name,
						})
						if err != nil {
							fmt.Println("error when try to get auto scaling application metrics during AutoScaleApplicationsAndNotifyScheduler :", err.Error())
							done <- true
							return
						}
						if len(metrics) == 0 {
							fmt.Println("No metrics found for auto scaling application", application.Name)
							done <- true
							return
						}

						// for _, metric := range metrics {
						// 	fmt.Println("Metrics of application", application.Name, ":", metric.String())
						// }

						var compareActualUsageToAcceptedPercentageResults []domain.CompareActualUsageToAcceptedPercentageResult
						acceptedCPUUsageThreshold := application.ScalabilitySpecifications.Data().CpuUsagePercentageThreshold
						acceptedMemoryUsageThreshold := application.ScalabilitySpecifications.Data().MemoryUsagePercentageThreshold
						for _, metric := range metrics {
							doesExceedCPUAcceptedPercentage, actualCPUUsage := domain.DoesUsageExceedsLimitAndHowMuchActually(metric.CPUUsage, metric.MaxCPUUsage, acceptedCPUUsageThreshold)
							doesExceedMemoryAcceptedPercentage, actualMemoryUsage := domain.DoesUsageExceedsLimitAndHowMuchActually(metric.MemoryUsage, metric.MaxMemoryUsage, acceptedMemoryUsageThreshold)
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
							if compareActualUsageToAcceptedPercentageResult.OneOfTheUsageExceedsAcceptedPercentage() {
								if application.ScalabilitySpecifications.Data().Replicas >= domain.MaxNumberOfReplicas {
									fmt.Println("Auto application", application.Name, "has reached the maximum number of replicas")

									if time.Since(lastNotifiedMaxReplicasDatetimeByApplication[application.ID]).Hours() < 4 {
										fmt.Println("Auto application", application.Name, "has exceeded accepted percentages, but notification email was sent less than 4 hours ago")
										return
									}
									success, err := scheduler.scalabilityNotificationService.SendApplicationCannotBeScaledUp(
										application.AdministratorEmail,
										application.Name,
										application.Namespace.Name,
										compareActualUsageToAcceptedPercentageResult,
										application.ContainerSpecifications,
										application.ScalabilitySpecifications.Data(),
									)
									if err != nil {
										fmt.Println("error when try to send application cannot be scaled up more notification mail :", err.Error())
										return
									}
									if success {
										fmt.Println("Application", application.Name, "cannot be scaled up more, email sent to", application.AdministratorEmail)

										lastNotifiedMaxReplicasDatetimeByApplication[application.ID] = time.Now()
									} else {
										fmt.Println("Application", application.Name, "cannot be scaled up more, but email not sent to", application.AdministratorEmail)
									}
									return
								}
								// 5. scale up/down the application if one of the usage exceeds the accepted percentage
								_, err := scheduler.scaleApplicationUseCase.Execute(application.ID, commands.UpdateApplication{}, applications.HorizontalUpScaling)
								if err != nil {
									fmt.Println("error when try to scale application during AutoScaleApplicationsAndNotifyScheduler :", err.Error())
									return
								}

								// 6. send email to the application administrator to notify him about the scaling
								success, err := scheduler.scalabilityNotificationService.SendApplicationScaledUp(
									application.AdministratorEmail,
									application.Name,
									application.Namespace.Name,
									compareActualUsageToAcceptedPercentageResult,
									application.ContainerSpecifications,
									application.ScalabilitySpecifications.Data(),
								)
								if err != nil {
									fmt.Println("Error when try to send application scalability notification mail during NotifyApplicationScalingRecommendationScheduler : " + err.Error())
									return
								}
								if success {
									fmt.Println("Auto application", application.Name, "has exceeded accepted percentages scaled up, email sent to", application.AdministratorEmail)
								} else {
									fmt.Println("Auto application", application.Name, "has exceeded accepted percentages scaled up, but email not sent to", application.AdministratorEmail)
								}
							}
							// else {
							// 	fmt.Println("Application", application.Name, "has not exceeded accepted percentages")
							// }
						}

						done <- true
					}(application)
				}
				for i := 0; i < routines; i++ {
					<-done
				}
			}
		}
	}()
}

func getAutoScaleAppSchedulerRepeatInterval() (int, error) {
	var repeatInterval int
	schedulerScaleApplicationAndNotifyInSeconds := os.Getenv("SCHEDULER_SCALE_APPLICATION_AND_NOTIFY_IN_SECONDS")
	if schedulerScaleApplicationAndNotifyInSeconds == "" {
		fmt.Println("SCHEDULER_RECOMMEND_APPLICATION_SCALING_IN_SECONDS is not set")
		return 0, fmt.Errorf("SCHEDULER_SCALE_APPLICATION_AND_NOTIFY_IN_SECONDS is not set")
	}
	repeatInterval, err := strconv.Atoi(schedulerScaleApplicationAndNotifyInSeconds)
	if err != nil {
		fmt.Println("Error when convert SCHEDULER_SCALE_APPLICATION_AND_NOTIFY_IN_SECONDS to int")
		return 0, fmt.Errorf("error when convert SCHEDULER_SCALE_APPLICATION_AND_NOTIFY_IN_SECONDS to int : %s", err.Error())
	}
	return repeatInterval, nil
}
