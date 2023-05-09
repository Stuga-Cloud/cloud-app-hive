package domain

// ApplicationScalabilitySpecifications is a struct that represents the scalability specifications of an application
type ApplicationScalabilitySpecifications struct {
	MinimumInstanceCount int32
	MaximumInstanceCount int32
	Replicas             int32
}
