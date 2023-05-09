package domain

// ApplicationType is an enum that represents the type of application : Load balanced, Single instance, etc.
type ApplicationType int

const (
	// SingleInstance is an application type that represents an application that is not load balanced
	SingleInstance ApplicationType = iota
	// LoadBalanced is an application type that represents an application that is load balanced
	LoadBalanced
)
