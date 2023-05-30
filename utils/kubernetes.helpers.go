package utils

import "fmt"

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
		"KB": "Ki",
		"MB": "Mi",
		"GB": "Gi",
		"TB": "Ti",
		"PB": "Pi",
		"EB": "Ei",
		"n":  "n",
		"m":  "m",
	}

	// Check if the unit exists in the map
	if k8sUnit, ok := realLifeUnitToK8sUnit[unit]; ok {
		return fmt.Sprintf("%s%s", value, k8sUnit)
	}

	// If unit is not found, return the original value
	return value
}

// DoesUsageExceedsLimitAndHowMuchActually function is used to check if the usage exceeds the limit and how much is the actual percentage of usage
func DoesUsageExceedsLimitAndHowMuchActually(usage string, limit string, acceptedUsagePercentage int64) (bool, float64) {
	usageNumeric := ConvertKubernetesResourceValueAndUnitToNumeric(usage)
	limitNumeric := ConvertKubernetesResourceValueAndUnitToNumeric(limit)

	if limitNumeric == 0 {
		return false, 0
	}

	actualPercentage := (usageNumeric / limitNumeric) * 100
	if actualPercentage > float64(acceptedUsagePercentage) {
		return true, actualPercentage
	}

	return false, actualPercentage
}

// Convert to numeric to be able to compare resources usage (for example 16Ki = 16, 16Mi = 16384)
func ConvertKubernetesResourceValueAndUnitToNumeric(value string) float64 {
	if len(value) == 0 {
		return 0
	}

	value, unit := getDigitsAndUnitFromString(value)

	// Map of Kubernetes units to real-life units
	realLifeUnitToK8sUnit := map[string]float64{
		"KB": 1,
		"MB": 1024,
		"GB": 1048576,
		"TB": 1073741824,
		"PB": 1099511627776,
		"EB": 1125899906842624,
		"n":  0.000000001,
		"m":  0.001,
	}

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
