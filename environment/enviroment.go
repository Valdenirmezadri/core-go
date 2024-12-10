package environment

type Environment string

func New(envString string) Environment {
	return Environment("").init(envString)
}

func (env Environment) init(envString string) Environment {
	if envString == "dev" {
		return Developer
	}

	return Production
}

const (
	Production Environment = "prod"
	Developer  Environment = "dev"
)
