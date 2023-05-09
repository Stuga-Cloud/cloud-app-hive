package domain

// DeployApplication is a command that represents the deployment of an application
type DeployApplication struct {
	Name                      string
	Image                     string
	Namespace                 string
	Port                      int
	ApplicationType           ApplicationType
	EnvironmentVariables      []ApplicationEnvironmentVariable
	Secrets                   []ApplicationSecret
	ContainerSpecifications   ApplicationContainerSpecifications
	ScalabilitySpecifications ApplicationScalabilitySpecifications
}
