package service

import "github.com/google/uuid"

// FIXME: Embed all values
type Config struct {
	ID         uuid.UUID `yaml:"id"`
	Settings   string    `yaml:"settings"`
	GatewayAPI struct {
		FallbackAddress string `yaml:"fallback_address"`
		Certificate     string `yaml:"certificate"`
	} `yaml:"gateway_api"`
}
