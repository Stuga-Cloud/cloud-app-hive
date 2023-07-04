package schedulers

import (
	"cloud-app-hive/domain"
	"cloud-app-hive/services"
	"cloud-app-hive/use_cases"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type NotifyAdminOnClusterExceededUsageScheduler struct {
	getClusterMetricsUseCase use_cases.GetClusterMetricsUseCase
	emailService             services.EmailService
}

func (scheduler NotifyAdminOnClusterExceededUsageScheduler) Launch() {
	fmt.Println("Starting 'NotifyAdminOnClusterExceededUsageScheduler' scheduler...")
	go func() {
		repeatInterval, err := getNotifyAdminOnClusterExceededUsageRepeatInterval()
		if err != nil {
			fmt.Println("Error when try to get notify admin on cluster exceeded usage scheduler repeat interval :", err.Error())
			return
		}
		ticker := time.NewTicker(time.Duration(repeatInterval) * time.Second)

		var lastNotificationSentToAdmin time.Time

		for {
			select {
			case <-ticker.C:
				clusterMetrics, err := scheduler.getClusterMetricsUseCase.Execute()
				if err != nil {
					fmt.Println("error when try to get cluster metrics during NotifyAdminOnClusterExceededUsageScheduler :", err.Error())
					continue
				}
				if clusterMetrics == nil {
					fmt.Println("No cluster metrics found")
					continue
				}

				notifyAdminWhenClusterNodesUsageIsAbovePercentageStr := os.Getenv("NOTIFY_ADMIN_WHEN_CLUSTER_NODES_USAGE_IS_ABOVE_PERCENTAGE")
				if notifyAdminWhenClusterNodesUsageIsAbovePercentageStr == "" {
					fmt.Println("NOTIFY_ADMIN_WHEN_CLUSTER_NODES_USAGE_IS_ABOVE_PERCENTAGE is not set")
					continue
				}
				notifyAdminWhenPercentageOfNodesExceedsUsageStr := os.Getenv("NOTIFY_ADMIN_WHEN_PERCENTAGE_OF_NODES_EXCEEDED_USAGE")
				if notifyAdminWhenPercentageOfNodesExceedsUsageStr == "" {
					fmt.Println("NOTIFY_ADMIN_WHEN_PERCENTAGE_OF_NODES_EXCEEDED_USAGE is not set")
					continue
				}

				notifyAdminWhenClusterNodesUsageIsAbovePercentage, err := strconv.ParseFloat(notifyAdminWhenClusterNodesUsageIsAbovePercentageStr, 64)
				if err != nil {
					fmt.Println("Error when convert NOTIFY_ADMIN_WHEN_CLUSTER_NODES_USAGE_IS_ABOVE_PERCENTAGE to float64 during NotifyAdminOnClusterExceededUsageScheduler : ", err.Error())
					continue
				}
				notifyAdminWhenPercentageOfNodesExceedsUsage, err := strconv.ParseFloat(notifyAdminWhenPercentageOfNodesExceedsUsageStr, 64)
				if err != nil {
					fmt.Printf("Error when convert NOTIFY_ADMIN_WHEN_PERCENTAGE_OF_NODES_EXCEEDED_USAGE to float64 during NotifyAdminOnClusterExceededUsageScheduler : %s\n", err.Error())
					continue
				}

				if domain.DoesPartOfNodesExceedCPUOrMemoryUsage(
					notifyAdminWhenClusterNodesUsageIsAbovePercentage,
					notifyAdminWhenPercentageOfNodesExceedsUsage,
					clusterMetrics.NodesComputedUsages,
				) {
					fmt.Println("Cluster is exceeding its limits !")
					clusterStateJSON, err := json.Marshal(clusterMetrics)
					if err != nil {
						fmt.Println("Error while marshalling cluster state: ", err)
						continue
					}
					fmt.Println("Cluster state: ", string(clusterStateJSON))
					if time.Since(lastNotificationSentToAdmin).Hours() > 4 {
						sendClusterExceededUsageEmailToAdmin(scheduler.emailService, clusterMetrics)
						lastNotificationSentToAdmin = time.Now()
					} else {
						fmt.Println("Email already sent to admin less than 8 hours ago")
					}
				}
			}
		}
	}()
}

// NodesMetrics
// NodesCapacities
func findCorrespondingNodeMetricsAndCapacities(
	nodeName string,
	clusterMetrics *domain.ClusterMetrics,
) (domain.NodeMetrics, domain.NodeCapacities, error) {
	for _, nodeMetrics := range clusterMetrics.NodesMetrics {
		if nodeMetrics.Name == nodeName {
			for _, nodeCapacities := range clusterMetrics.NodesCapacities {
				if nodeCapacities.Name == nodeName {
					return nodeMetrics, nodeCapacities, nil
				}
			}
		}
	}

	return domain.NodeMetrics{}, domain.NodeCapacities{}, fmt.Errorf("no corresponding node metrics and capacities found for node %s", nodeName)
}

func sendClusterExceededUsageEmailToAdmin(emailService services.EmailService, clusterMetrics *domain.ClusterMetrics) error {
	subject := "Cluster is exceeding its limits"
	body := "Cluster is exceeding its limits. Here are cluster metrics : \n"
	for _, nodeComputedUsage := range clusterMetrics.NodesComputedUsages {
		nodeMetrics, nodeCapacities, err := findCorrespondingNodeMetricsAndCapacities(
			nodeComputedUsage.Name,
			clusterMetrics,
		)
		if err != nil {
			fmt.Println("Error when try to find corresponding node metrics and capacities :", err.Error())
			return fmt.Errorf("error when try to find corresponding node metrics and capacities : %s", err.Error())
		}

		body += fmt.Sprintf(
			"Node %s\n",
			nodeComputedUsage.Name,
		)
		body += fmt.Sprintf(
			`
			- CPU = using %2.f%% of total CPU - absolute value : %v / capacity : %v CPU\n
			- Memory = using %2.f%% of total memory - absolute value : %v / capacity : %v\n
			`,
			nodeComputedUsage.CPUUsageInPercentage,
			nodeMetrics.ReadableCPUUsage,
			nodeCapacities.ReadableCPU,
			nodeComputedUsage.MemoryUsageInPercentage,
			nodeMetrics.ReadableMemoryUsage,
			nodeCapacities.ReadableMemory,
		)
	}

	htmlBody :=
		`
		<p>Cluster is exceeding its limits. Here are cluster metrics : </p>
		<ul>
		`

	for _, nodeComputedUsage := range clusterMetrics.NodesComputedUsages {
		nodeMetrics, nodeCapacities, err := findCorrespondingNodeMetricsAndCapacities(
			nodeComputedUsage.Name,
			clusterMetrics,
		)
		if err != nil {
			fmt.Println("Error when try to find corresponding node metrics and capacities :", err.Error())
			return fmt.Errorf("error when try to find corresponding node metrics and capacities : %s", err.Error())
		}

		htmlBody += fmt.Sprintf(
			`
			<li>Node %s</li>
			<ul>
				<li>CPU = using %2.f%% of total CPU - absolute value : %v / capacity : %v</li>
				<li>Memory = using %2.f%% of total memory - absolute value : %v / capacity : %v</li>
			</ul>
			`,
			nodeComputedUsage.Name,
			nodeComputedUsage.CPUUsageInPercentage,
			nodeMetrics.ReadableCPUUsage,
			nodeCapacities.ReadableCPU,
			nodeComputedUsage.MemoryUsageInPercentage,
			nodeMetrics.ReadableMemoryUsage,
			nodeCapacities.ReadableMemory,
		)
	}
	htmlBody += "</ul>"

	adminEmail := os.Getenv("ADMIN_EMAIL")
	if adminEmail == "" {
		fmt.Println("ADMIN_EMAIL is not set")
		return fmt.Errorf("ADMIN_EMAIL is not set")
	}
	copyCarbonEmails := os.Getenv("ADMIN_EMAIL_COPY_CARBON")
	if copyCarbonEmails == "" {
		fmt.Println("ADMIN_EMAIL_COPY_CARBON is not set")
		return fmt.Errorf("ADMIN_EMAIL_COPY_CARBON is not set")
	}
	copyCarbonCopyEmails := strings.Split(copyCarbonEmails, ",")

	fmt.Printf("Sending %s email to admin %s and copy carbon copy emails %v\n", subject, adminEmail, copyCarbonCopyEmails)
	err := emailService.Send(
		adminEmail,
		subject,
		body,
		htmlBody,
		copyCarbonCopyEmails,
	)
	if err != nil {
		fmt.Println("error when try to send application scalability notification mail during NotifyApplicationScalingRecommendationScheduler :", err.Error())
		return fmt.Errorf("error when try to send application scalability notification mail during NotifyApplicationScalingRecommendationScheduler : %s", err.Error())
	}

	return nil
}

func getNotifyAdminOnClusterExceededUsageRepeatInterval() (int, error) {
	var repeatInterval int
	SchedulerRecommendApplicationScalingInSeconds := os.Getenv("SCHEDULER_NOTIFY_ADMIN_ON_CLUSTER_EXCEEDED_USAGE_IN_SECONDS")
	if SchedulerRecommendApplicationScalingInSeconds == "" {
		fmt.Println("SCHEDULER_NOTIFY_ADMIN_ON_CLUSTER_EXCEEDED_USAGE_IN_SECONDS is not set")
		return 0, fmt.Errorf("SCHEDULER_NOTIFY_ADMIN_ON_CLUSTER_EXCEEDED_USAGE_IN_SECONDS is not set")
	}
	repeatInterval, err := strconv.Atoi(SchedulerRecommendApplicationScalingInSeconds)
	if err != nil {
		fmt.Println("Error when convert SCHEDULER_NOTIFY_ADMIN_ON_CLUSTER_EXCEEDED_USAGE_IN_SECONDS to int")
		return 0, fmt.Errorf("error when convert SCHEDULER_NOTIFY_ADMIN_ON_CLUSTER_EXCEEDED_USAGE_IN_SECONDS to int during NotifyApplicationScalingRecommendationScheduler : %s", err.Error())
	}
	return repeatInterval, nil
}
