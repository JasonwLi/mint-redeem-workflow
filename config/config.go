package config

type ServiceConfig struct {
	BraleAuth string
}

func NewServiceConfig() (*ServiceConfig, error) {
	jwt := "test"

	return &ServiceConfig{
		BraleAuth: jwt,
	}, nil
}
