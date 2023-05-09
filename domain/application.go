package domain

// Application is a struct that represents a user's application
type Application struct {
	Name                      string
	Desc                      string
	Image                     string
	Namespace                 string
	Port                      int
	Zone                      string // The zone where the application is deployed (e.g. eu-west-1)
	ApplicationType           ApplicationType
	EnvironmentVariables      []ApplicationEnvironmentVariable
	Secrets                   []ApplicationSecret
	ContainerSpecifications   ApplicationContainerSpecifications
	ScalabilitySpecifications ApplicationScalabilitySpecifications
}
