package domain

import (
	"fmt"
	"math"
	"strconv"
)

func getDigitsAndUnitFromString(value string) (string, string) {
	digits := ""
	unit := ""
	for _, char := range value {
		if char >= '0' && char <= '9' {
			digits += string(char)
		} else {
			unit += string(char)
		}
	}
	return digits, unit
}

// ConvertK8sResourceToReadableHumanValueAndUnit function is used to convert k8s resource values and units to real-life values and units
func ConvertK8sResourceToReadableHumanValueAndUnit(k8sValue string) string {
	if len(k8sValue) == 0 {
		return "0"
	}

	value, unit := getDigitsAndUnitFromString(k8sValue)

	// Map of Kubernetes units to real-life units
	k8sUnitToRealLifeUnit := map[string]string{
		"Ki": "KB",
		"Mi": "MB",
		"Gi": "GB",
		"Ti": "TB",
		"Pi": "PB",
		"Ei": "EB",
		"n":  "n",
		"m":  "m",
	}

	// Check if the unit exists in the map
	if realLifeUnit, ok := k8sUnitToRealLifeUnit[unit]; ok {
		return fmt.Sprintf("%s%s", value, realLifeUnit)
	}

	// If unit is not found, return the original value
	return k8sValue
}

func ConvertReadableHumanValueAndUnitToK8sResource(value string) string {
	if len(value) == 0 {
		return "0"
	}

	value, unit := getDigitsAndUnitFromString(value)
	println("Value & unit: ", value, unit)

	// Map of Kubernetes units to real-life units
	realLifeUnitToK8sUnit := map[string]string{
		"KB":   "Ki",
		"MB":   "Mi",
		"GB":   "Gi",
		"TB":   "Ti",
		"PB":   "Pi",
		"EB":   "Ei",
		"nCPU": "n",
		"mCPU": "m",
	}

	// Check if the unit exists in the map
	if k8sUnit, ok := realLifeUnitToK8sUnit[unit]; ok {
		return fmt.Sprintf("%s%s", value, k8sUnit)
	}

	fmt.Println("Unit not found: ", unit)

	// If unit is not found, return the original value
	return value
}

// DoesUsageExceedsLimitAndHowMuchActually function is used to check if the usage exceeds the limit and how much is the actual percentage of usage
func DoesUsageExceedsLimitAndHowMuchActually(usage string, limit string, acceptedUsagePercentage float64) (bool, float64) {
	usageNumeric := ConvertKubernetesResourceValueAndUnitToNumeric(usage)
	limitNumeric := ConvertKubernetesResourceValueAndUnitToNumeric(limit)

	if limitNumeric == 0 {
		return false, 0
	}

	actualPercentage := (usageNumeric / limitNumeric) * 100
	if actualPercentage > acceptedUsagePercentage {
		return true, actualPercentage
	}

	return false, actualPercentage
}

// Map of Kubernetes units to real-life units
const oneKiloByte = 1024

var realLifeUnitToK8sUnit = map[string]float64{
	"KB": math.Pow(oneKiloByte, 1),
	"MB": math.Pow(oneKiloByte, 2),
	"GB": math.Pow(oneKiloByte, 3),
	"TB": math.Pow(oneKiloByte, 4),
	"PB": math.Pow(oneKiloByte, 5),
	"EB": math.Pow(oneKiloByte, 6),
	"n":  math.Pow(10, -9),
	"u":  math.Pow(10, -6),
	"m":  math.Pow(10, -3),
}

// Convert to numeric to be able to compare resources usage (for example 16Ki = 16, 16Mi = 16384)
func ConvertKubernetesResourceValueAndUnitToNumeric(value string) float64 {
	if len(value) == 0 {
		return 0
	}

	value, unit := getDigitsAndUnitFromString(value)

	// Check if the unit exists in the map
	if k8sUnit, ok := realLifeUnitToK8sUnit[unit]; ok {
		return k8sUnit * ConvertStringToFloat64(value)
	}

	// If unit is not found, return the original value
	return 0
}

func ConvertStringToFloat64(value string) float64 {
	if len(value) == 0 {
		return 0
	}

	var result float64
	_, err := fmt.Sscanf(value, "%f", &result)
	if err != nil {
		return 0
	}

	return result
}

func DivideFloat64s(a float64, b float64) float64 {
	if b == 0 {
		return 0
	}

	return a / b
}

func ComputeNodesUsagesFromMetricsAndCapacities(nodeMetrics []NodeMetrics, nodeCapacities []NodeCapacities) ([]NodeComputedUsage, error) {
	if len(nodeMetrics) != len(nodeCapacities) {
		fmt.Println("Mismatched number of metrics and capacities : ", len(nodeMetrics), len(nodeCapacities))
		return nil, fmt.Errorf("mismatched number of metrics and capacities")
	}

	var nodeComputedUsages []NodeComputedUsage

	for _, nodeMetric := range nodeMetrics {
		for _, nodeCapacity := range nodeCapacities {
			if nodeMetric.Name == nodeCapacity.Name {
				cpuUsage := ConvertKubernetesResourceValueAndUnitToNumeric(nodeMetric.ReadableCPUUsage)
				cpuLimitConverted, err := strconv.Atoi(nodeCapacity.CPULimit)
				if err != nil {
					fmt.Println("Error while converting CPU limit to int: ", err)
				}
				cpuCapacity := float64(cpuLimitConverted)
				cpuUsagePercentage := DivideFloat64s(cpuUsage, cpuCapacity)

				memoryUsage := ConvertKubernetesResourceValueAndUnitToNumeric(nodeMetric.ReadableMemoryUsage)
				memoryCapacity := ConvertKubernetesResourceValueAndUnitToNumeric(nodeCapacity.ReadableMemory)
				memoryUsagePercentage := DivideFloat64s(memoryUsage, memoryCapacity)

				ephemeralStorageUsage := ConvertKubernetesResourceValueAndUnitToNumeric(nodeMetric.ReadableEphemeralStorageUsage)
				ephemeralStorageCapacity := ConvertKubernetesResourceValueAndUnitToNumeric(nodeCapacity.ReadableEphemeralStorage)
				ephemeralStorageUsagePercentage := DivideFloat64s(ephemeralStorageUsage, ephemeralStorageCapacity)

				storageUsage := ConvertKubernetesResourceValueAndUnitToNumeric(nodeMetric.ReadableStorageUsage)
				storageCapacity := ConvertKubernetesResourceValueAndUnitToNumeric(nodeCapacity.ReadableStorage)
				storagesUsagePercentage := DivideFloat64s(storageUsage, storageCapacity)

				// fmt.Println("Node: ", nodeMetric.Name, " CPU usage: ", cpuUsagePercentage, " Memory usage: ", memoryUsagePercentage, " Storage usage: ", storagesUsagePercentage, " Ephemeral storage usage: ", ephemeralStorageUsagePercentage)

				nodeComputedUsages = append(nodeComputedUsages, NodeComputedUsage{
					Name:                              nodeMetric.Name,
					CPUUsageInPercentage:              cpuUsagePercentage * 100,
					MemoryUsageInPercentage:           memoryUsagePercentage * 100,
					StorageUsageInPercentage:          storagesUsagePercentage * 100,
					EphemeralStorageUsageInPercentage: ephemeralStorageUsagePercentage * 100,
				})
			}
		}
	}

	return nodeComputedUsages, nil
}

func doesHalfNodesExceedCPUOrMemoryUsage(nodeMetrics []NodeMetrics, nodeCapacities []NodeCapacities) bool {
	if len(nodeMetrics) != len(nodeCapacities) {
		fmt.Println("Mismatched number of metrics and capacities.")
		return false
	}

	threshold := 0.8
	exceedingNodes := 0

	for i := 0; i < len(nodeMetrics); i++ {
		cpuUsage, _ := strconv.Atoi(nodeMetrics[i].CPUUsage)
		cpuLimit, _ := strconv.Atoi(nodeCapacities[i].CPULimit)
		memoryUsage, _ := strconv.Atoi(nodeMetrics[i].MemoryUsage)
		memoryLimit, _ := strconv.Atoi(nodeCapacities[i].MemoryLimit)

		cpuUsagePercentage := float64(cpuUsage) / float64(cpuLimit)
		memoryUsagePercentage := float64(memoryUsage) / float64(memoryLimit)

		if cpuUsagePercentage >= threshold || memoryUsagePercentage >= threshold {
			exceedingNodes++
		}
	}

	halfNodes := len(nodeMetrics) / 2

	if exceedingNodes >= halfNodes {
		return true
	}

	return false
}

func DoesPartOfNodesExceedCPUOrMemoryUsage(
	acceptedUsagePercentage float64, // 50 => 50%
	nodesPercentageThatExceed float64, // 50 => 50%
	nodesComputedUsages []NodeComputedUsage,
) bool {
	usageThreshold := acceptedUsagePercentage
	exceedingNodes := 0

	for _, nodeComputedUsage := range nodesComputedUsages {
		if nodeComputedUsage.CPUUsageInPercentage >= usageThreshold || nodeComputedUsage.MemoryUsageInPercentage >= usageThreshold {
			exceedingNodes++
		}
	}

	partOfNodes := float64(len(nodesComputedUsages)) * (nodesPercentageThatExceed / 100)
	if partOfNodes < 1 {
		partOfNodes = 1
	}
	if float64(exceedingNodes) >= math.Floor(partOfNodes) {
		return true
	}

	return false
}
