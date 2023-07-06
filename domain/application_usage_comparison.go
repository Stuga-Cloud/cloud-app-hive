package domain

type UsageComparisonResult struct {
	DoesExceedAcceptedPercentage bool
	ActualUsage                  float64
}

type CompareActualUsageToAcceptedPercentageResult struct {
	PodName                     string
	CPUUsageResult              UsageComparisonResult
	MemoryUsageResult           UsageComparisonResult
	EphemeralStorageUsageResult UsageComparisonResult
}

func (r *CompareActualUsageToAcceptedPercentageResult) OneOfTheUsageExceedsAcceptedPercentage() bool {
	return r.CPUUsageResult.DoesExceedAcceptedPercentage || r.MemoryUsageResult.DoesExceedAcceptedPercentage || r.EphemeralStorageUsageResult.DoesExceedAcceptedPercentage
}

func (r *CompareActualUsageToAcceptedPercentageResult) CPUAndMemoryUsageExceedsAcceptedPercentage() bool {
	return r.CPUUsageResult.DoesExceedAcceptedPercentage && r.MemoryUsageResult.DoesExceedAcceptedPercentage
}
